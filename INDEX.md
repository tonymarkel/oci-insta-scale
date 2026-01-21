# 📖 Documentation Index

Welcome to OCI Insta-Scale! This index will help you find the right documentation for your needs.

## 🚀 Getting Started (Start Here!)

| Document | Purpose | When to Read |
|----------|---------|--------------|
| [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) | **Overview of everything** | Read this FIRST |
| [CHECKLIST.md](CHECKLIST.md) | **Step-by-step setup guide** | Follow during setup |
| [QUICKSTART.md](QUICKSTART.md) | **Quick setup & commands** | Reference while setting up |

## 📚 Detailed Documentation

| Document | Purpose | When to Read |
|----------|---------|--------------|
| [README.md](README.md) | **Complete feature docs** | When you need details |
| [EXAMPLES.md](EXAMPLES.md) | **Configuration examples** | When configuring |
| [ARCHITECTURE.md](ARCHITECTURE.md) | **Project structure** | Understanding the codebase |
| [WORKFLOW.md](WORKFLOW.md) | **Visual diagrams** | Understanding workflows |

## 🗂️ Configuration Files

| File | Purpose | Notes |
|------|---------|-------|
| `config.yaml` | **Your actual config** | ⚠️ Gitignored - contains secrets |
| `config.example.yaml` | **Template config** | Copy this to create config.yaml |

## 💻 Source Code

| File | Purpose | Language |
|------|---------|----------|
| `main.go` | Instance creation | Go |
| `capacity-manager.go` | Reservation manager | Go |
| `manage-instances.sh` | Batch operations | Bash |
| `Makefile` | Build automation | Make |

## 🎯 Quick Navigation

### I want to...

**→ Understand what this project does**
- Start with: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

**→ Set up the project for the first time**
1. [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Overview
2. [CHECKLIST.md](CHECKLIST.md) - Follow step-by-step
3. [QUICKSTART.md](QUICKSTART.md) - Reference commands

**→ Configure my instances**
- Check: [EXAMPLES.md](EXAMPLES.md) - Multiple configuration examples

**→ Understand how it works**
- Read: [WORKFLOW.md](WORKFLOW.md) - Visual diagrams
- Then: [ARCHITECTURE.md](ARCHITECTURE.md) - Project structure

**→ Find a specific command**
- Quick reference: [QUICKSTART.md](QUICKSTART.md)
- Full details: [README.md](README.md)

**→ Troubleshoot an issue**
1. [CHECKLIST.md](CHECKLIST.md) - Common Issues section
2. [README.md](README.md) - Troubleshooting section
3. [QUICKSTART.md](QUICKSTART.md) - Error handling

**→ Learn about capacity reservations**
- [README.md](README.md) - Capacity Reservations section
- [EXAMPLES.md](EXAMPLES.md) - Example 2 & 4

**→ Manage existing instances**
- [QUICKSTART.md](QUICKSTART.md) - Section 7: Cleanup
- [README.md](README.md) - Cleanup section

**→ Understand the code structure**
- [ARCHITECTURE.md](ARCHITECTURE.md)
- [WORKFLOW.md](WORKFLOW.md)

## 📋 Documentation by Role

### 👤 First-Time User
1. [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
2. [CHECKLIST.md](CHECKLIST.md)
3. [EXAMPLES.md](EXAMPLES.md) - Example 1

### 👔 Operations Engineer
1. [QUICKSTART.md](QUICKSTART.md)
2. [README.md](README.md)
3. [EXAMPLES.md](EXAMPLES.md) - Example 2

### 👨‍💻 Developer
1. [ARCHITECTURE.md](ARCHITECTURE.md)
2. [WORKFLOW.md](WORKFLOW.md)
3. Source files: `main.go`, `capacity-manager.go`

### 🔬 Researcher/HPC User
1. [EXAMPLES.md](EXAMPLES.md) - Example 4
2. [README.md](README.md) - Capacity Reservations

### 💰 Cost-Conscious User
1. [EXAMPLES.md](EXAMPLES.md) - Example 3 (minimal) & 5 (ARM)
2. [QUICKSTART.md](QUICKSTART.md) - Cleanup section

## 📊 Documentation Stats

| Document | Lines | Focus |
|----------|-------|-------|
| PROJECT_SUMMARY.md | ~300 | Complete overview |
| README.md | ~250 | Feature documentation |
| QUICKSTART.md | ~200 | Quick reference |
| CHECKLIST.md | ~250 | Setup guide |
| EXAMPLES.md | ~400 | Config examples |
| WORKFLOW.md | ~300 | Visual workflows |
| ARCHITECTURE.md | ~250 | Project structure |

## 🔍 Search Guide

Looking for something specific? Use these search terms:

**Commands & Usage:**
- Search: "make", "run", "./oci-insta-scale"
- Files: [QUICKSTART.md](QUICKSTART.md), [README.md](README.md)

**Configuration:**
- Search: "config.yaml", "ocid", "shape"
- Files: [EXAMPLES.md](EXAMPLES.md), [README.md](README.md)

**Capacity Reservations:**
- Search: "capacity-manager", "reservation"
- Files: [README.md](README.md), [EXAMPLES.md](EXAMPLES.md)

**Troubleshooting:**
- Search: "error", "troubleshoot", "issue"
- Files: [CHECKLIST.md](CHECKLIST.md), [README.md](README.md)

**Architecture:**
- Search: "workflow", "diagram", "structure"
- Files: [WORKFLOW.md](WORKFLOW.md), [ARCHITECTURE.md](ARCHITECTURE.md)

## 🆘 Help Decision Tree

```
Having an issue?
│
├─ Don't know where to start?
│  └─ Read: PROJECT_SUMMARY.md
│
├─ Setup not working?
│  └─ Follow: CHECKLIST.md step-by-step
│
├─ Configuration question?
│  └─ Check: EXAMPLES.md for similar use case
│
├─ Command not found?
│  └─ Reference: QUICKSTART.md or Makefile
│
├─ Error message?
│  └─ Check: CHECKLIST.md Common Issues
│  └─ Then: README.md Troubleshooting
│
└─ Want to understand internals?
   └─ Read: WORKFLOW.md then ARCHITECTURE.md
```

## ⚡ Quick Command Reference

```bash
# Documentation
cat PROJECT_SUMMARY.md  # Start here
cat CHECKLIST.md        # Setup guide
make help               # Available commands

# Setup
make setup-config       # Create config from template
make deps               # Install dependencies
make build              # Build binaries

# Running
make dry-run           # Test configuration
make run               # Create instances
make list-reservations # List capacity reservations

# Managing
./manage-instances.sh list -c <id>      # List instances
./manage-instances.sh status -c <id>    # Check status
./manage-instances.sh terminate -c <id> # Clean up
```

## 📦 Complete File Tree

```
oci-insta-scale/
├── 📚 Documentation (Read These!)
│   ├── INDEX.md               ← You are here
│   ├── PROJECT_SUMMARY.md     ← Start here!
│   ├── CHECKLIST.md           ← Setup guide
│   ├── QUICKSTART.md          ← Quick reference
│   ├── README.md              ← Full documentation
│   ├── EXAMPLES.md            ← Config examples
│   ├── WORKFLOW.md            ← Diagrams
│   └── ARCHITECTURE.md        ← Structure
│
├── 💻 Source Code
│   ├── main.go                # Instance creator
│   ├── capacity-manager.go    # Reservation manager
│   ├── manage-instances.sh    # Batch operations
│   └── Makefile               # Build automation
│
├── ⚙️  Configuration
│   ├── config.yaml            # Your config (gitignored)
│   ├── config.example.yaml    # Template
│   ├── go.mod                 # Go dependencies
│   └── go.sum                 # Checksums
│
└── 🔨 Build Artifacts
    ├── oci-insta-scale        # Main binary
    └── capacity-manager       # Manager binary
```

## 🎓 Learning Path

### Beginner (Never used OCI)
1. Read [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
2. Follow [CHECKLIST.md](CHECKLIST.md) completely
3. Try [EXAMPLES.md](EXAMPLES.md) - Example 1 (basic)
4. Refer to [QUICKSTART.md](QUICKSTART.md) as needed

### Intermediate (Used OCI, new to this tool)
1. Skim [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
2. Read [QUICKSTART.md](QUICKSTART.md)
3. Pick example from [EXAMPLES.md](EXAMPLES.md)
4. Run and iterate

### Advanced (Building/Modifying)
1. Review [ARCHITECTURE.md](ARCHITECTURE.md)
2. Study [WORKFLOW.md](WORKFLOW.md)
3. Read source: `main.go`, `capacity-manager.go`
4. Refer to [README.md](README.md) for API details

## 🔖 Bookmarks

Keep these handy:

- **Daily use**: [QUICKSTART.md](QUICKSTART.md)
- **Configuration**: [EXAMPLES.md](EXAMPLES.md)
- **Troubleshooting**: [CHECKLIST.md](CHECKLIST.md)
- **Reference**: [README.md](README.md)

## 📝 Notes

- All documentation is in Markdown format
- Files prefixed with capitals are documentation
- Use `grep` to search across all docs:
  ```bash
  grep -r "search term" *.md
  ```
- Keep this INDEX.md open in a tab for quick reference

---

**Happy scaling with OCI! 🚀**
