# OCI Insta-Scale Workflow

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    OCI Insta-Scale                          │
│                                                             │
│  ┌──────────────┐    ┌──────────────┐   ┌───────────────┐ │
│  │   main.go    │    │  capacity-   │   │   manage-     │ │
│  │              │    │  manager.go  │   │ instances.sh  │ │
│  │  Instance    │    │              │   │               │ │
│  │  Creation    │    │ Reservation  │   │   Batch Ops   │ │
│  └──────┬───────┘    └──────┬───────┘   └───────┬───────┘ │
│         │                   │                    │         │
└─────────┼───────────────────┼────────────────────┼─────────┘
          │                   │                    │
          │                   │                    │
          └───────────────────┴────────────────────┘
                              │
                    ┌─────────▼─────────┐
                    │  OCI Go SDK v65   │
                    └─────────┬─────────┘
                              │
                    ┌─────────▼─────────┐
                    │   OCI REST API    │
                    └─────────┬─────────┘
                              │
          ┌───────────────────┴───────────────────┐
          │                                       │
    ┌─────▼─────┐                        ┌───────▼────────┐
    │  Capacity │                        │   Compute      │
    │Reservations│                        │  Instances     │
    └───────────┘                        └────────────────┘
```

## Data Flow

### 1. Configuration Loading
```
config.yaml
    │
    ├─ OCI Credentials (tenancy, user, key, region)
    ├─ Compartment ID
    └─ Instance Settings
        ├─ Shape & Resources (OCPUs, Memory)
        ├─ Image & Network (Image OCID, Subnet)
        ├─ Capacity Reservation (Optional)
        └─ Tags & Metadata
```

### 2. Instance Creation Flow

```
Start
  │
  ▼
Load config.yaml
  │
  ▼
Validate Configuration
  │ ├─ Missing fields? ─> ERROR
  ▼
Create OCI Client
  │
  ▼
Dry Run? ──Yes──> Display Config & Exit
  │
  No
  ▼
Loop: For each instance (1 to N)
  │
  ├─> Build Launch Request
  │   ├─ Shape Configuration
  │   ├─ Image Details
  │   ├─ Network Details
  │   ├─ Capacity Reservation (if set)
  │   └─ Metadata & Tags
  │
  ├─> Call OCI API: LaunchInstance
  │   │
  │   ├─ Success ──> Store Instance OCID
  │   │              │
  │   │              ▼
  │   │          Log Success
  │   │              │
  │   │              ▼
  │   │          Sleep 1s (rate limit)
  │   │
  │   └─ Failure ──> Log Error, Continue
  │
  ▼
End Loop
  │
  ▼
Print Summary
  │
  ▼
Wait for RUNNING state? ──No──> Exit
  │
  Yes
  ▼
Loop: For each instance
  │
  ├─> Poll Instance State
  │   │
  │   ├─ RUNNING ──> Log Success
  │   │
  │   ├─ PROVISIONING/STARTING ──> Wait 10s, Retry
  │   │
  │   └─ TERMINATED ──> Log Error
  │
  ▼
End Loop
  │
  ▼
Done!
```

### 3. Capacity Reservation Flow

```
┌─────────────────────────────────────────────────────────┐
│         Capacity Reservation Management                 │
└─────────────────────────────────────────────────────────┘

CREATE RESERVATION
─────────────────
capacity-manager -create
  │
  ├─ Inputs: Name, AD, Shape, Count, OCPUs, Memory
  │
  ├─> Build Reservation Config
  │   └─ Instance Shape Config (for Flex shapes)
  │
  ├─> Call OCI API: CreateComputeCapacityReservation
  │
  ├─> Return: Reservation OCID
  │
  └─> User adds OCID to config.yaml


USE RESERVATION
───────────────
oci-insta-scale (with capacity_reservation_id set)
  │
  ├─> Read capacity_reservation_id from config
  │
  ├─> Include in LaunchInstance request
  │
  └─> Instances use reserved capacity


LIST RESERVATIONS
─────────────────
capacity-manager -list
  │
  ├─> Call OCI API: ListComputeCapacityReservations
  │
  └─> Display: Name, OCID, State, AD, Usage


DELETE RESERVATION
──────────────────
capacity-manager -delete -id <ocid>
  │
  ├─> Call OCI API: DeleteComputeCapacityReservation
  │
  └─> Confirm deletion
```

## Component Interactions

```
┌────────────────┐
│     User       │
└───────┬────────┘
        │
        │ (1) Configure
        ▼
┌────────────────┐
│  config.yaml   │◄────────────────┐
└───────┬────────┘                 │
        │                          │
        │ (2) Run                  │
        ▼                          │
┌────────────────┐                 │
│ oci-insta-     │                 │ (6) Update with
│    scale       │                 │ Reservation ID
└───────┬────────┘                 │
        │                          │
        │ (3) Create               │
        │ Instances                │
        ▼                          │
┌────────────────┐                 │
│   OCI API      │                 │
└───────┬────────┘                 │
        │                          │
        │ (4) Optionally           │
        │ Create Reservation       │
        │                          │
        ▼                          │
┌────────────────┐                 │
│   capacity-    │─────────────────┘
│   manager      │ (5) Return OCID
└───────┬────────┘
        │
        │ (7) Manage
        │ Lifecycle
        ▼
┌────────────────┐
│   manage-      │
│ instances.sh   │
└────────────────┘
```

## State Diagram: Instance Lifecycle

```
┌─────────┐
│ Request │
│ Create  │
└────┬────┘
     │
     ▼
┌─────────────┐
│ PROVISIONING│
└────┬────────┘
     │
     ▼
┌──────────┐
│ STARTING │
└────┬─────┘
     │
     ▼
┌──────────┐     Stop      ┌──────────┐
│ RUNNING  │◄──────────────│ STOPPED  │
└────┬─────┘               └────┬─────┘
     │                          │
     │ Start                    │
     └──────────────────────────┘
     │
     │ Terminate
     ▼
┌──────────────┐
│ TERMINATING  │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ TERMINATED   │
└──────────────┘
```

## Quick Command Reference

```
┌──────────────────────────────────────────────────────┐
│              COMMON WORKFLOWS                        │
└──────────────────────────────────────────────────────┘

1. FIRST TIME SETUP
   ─────────────────
   $ make setup-config
   $ vi config.yaml         # Edit with your settings
   $ make deps
   $ make build

2. CREATE INSTANCES (No Reservation)
   ──────────────────────────────────
   $ make dry-run          # Test first
   $ make run              # Create instances

3. CREATE WITH RESERVATION
   ────────────────────────
   $ make list-reservations    # Check existing
   $ ./capacity-manager -create \
       -name "prod" -ad "..." -shape "..." -count 10
   $ vi config.yaml            # Add reservation OCID
   $ make run

4. MANAGE INSTANCES
   ─────────────────
   $ ./manage-instances.sh list -c <compartment-id>
   $ ./manage-instances.sh stop -c <id> -p "test-"
   $ ./manage-instances.sh start -c <id> -p "test-"
   $ ./manage-instances.sh status -c <id>

5. CLEANUP
   ────────
   $ ./manage-instances.sh terminate -c <id> -p "test-"
   $ ./capacity-manager -delete -id <reservation-ocid>

┌──────────────────────────────────────────────────────┐
│              CAPACITY PLANNING                       │
└──────────────────────────────────────────────────────┘

Shape: VM.Standard.E4.Flex
├─ 1 OCPU + 6 GB RAM
├─ 10 instances = 10 OCPUs, 60 GB RAM
└─ Cost: ~$0.06/hr per instance

Reservation Benefits:
├─ Guaranteed capacity
├─ Reserved for your tenancy
├─ No additional cost
└─ Can be released when not needed
```

## Error Handling

```
Error Type              Handling Strategy
──────────              ─────────────────
Out of Capacity      -> Create Capacity Reservation
Authentication       -> Check credentials & API keys
Shape Unavailable    -> Try different AD or shape
Subnet Not Found     -> Verify subnet OCID
Rate Limiting        -> Built-in 1s delay between requests
Instance Launch Fail -> Continue with remaining instances
Network Error        -> Retry with exponential backoff
```

## Best Practices

```
┌─────────────────────────────────────────────┐
│ 1. ALWAYS DRY-RUN FIRST                     │
│    $ make dry-run                           │
├─────────────────────────────────────────────┤
│ 2. USE CAPACITY RESERVATIONS                │
│    For production workloads                 │
├─────────────────────────────────────────────┤
│ 3. TAG YOUR RESOURCES                       │
│    Cost tracking & organization             │
├─────────────────────────────────────────────┤
│ 4. TEST IN DEV FIRST                        │
│    Separate configs per environment         │
├─────────────────────────────────────────────┤
│ 5. CLEAN UP REGULARLY                       │
│    Terminate unused instances               │
└─────────────────────────────────────────────┘
```
