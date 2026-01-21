#!/bin/bash

# Batch Instance Management Script
# Helps manage multiple instances created by oci-insta-scale

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

function print_usage() {
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  list                    List all instances matching prefix"
    echo "  stop                    Stop all instances matching prefix"
    echo "  start                   Start all instances matching prefix"
    echo "  terminate               Terminate all instances matching prefix"
    echo "  status                  Show detailed status of instances"
    echo ""
    echo "Options:"
    echo "  -c, --compartment-id    Compartment OCID (required)"
    echo "  -p, --prefix            Instance name prefix (default: 'instance-')"
    echo "  -y, --yes               Auto-confirm destructive actions"
    echo ""
    echo "Examples:"
    echo "  $0 list -c ocid1.compartment..."
    echo "  $0 stop -c ocid1.compartment... -p 'prod-'"
    echo "  $0 terminate -c ocid1.compartment... -p 'test-' -y"
}

function check_oci_cli() {
    if ! command -v oci &> /dev/null; then
        echo -e "${RED}Error: OCI CLI is not installed${NC}"
        echo "Install it from: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm"
        exit 1
    fi
}

function list_instances() {
    echo -e "${YELLOW}Listing instances with prefix: $PREFIX${NC}"
    oci compute instance list \
        --compartment-id "$COMPARTMENT_ID" \
        --display-name "$PREFIX*" \
        --output table \
        --query 'data[*].{Name:"display-name",OCID:id,State:"lifecycle-state",AD:"availability-domain",Created:"time-created"}'
}

function get_instance_ids() {
    oci compute instance list \
        --compartment-id "$COMPARTMENT_ID" \
        --display-name "$PREFIX*" \
        --query 'data[].id' \
        --raw-output | tr -d '[]", '
}

function stop_instances() {
    local instance_ids=$(get_instance_ids)
    
    if [ -z "$instance_ids" ]; then
        echo -e "${YELLOW}No instances found with prefix: $PREFIX${NC}"
        return
    fi
    
    echo -e "${YELLOW}Stopping instances with prefix: $PREFIX${NC}"
    
    if [ "$AUTO_CONFIRM" != "true" ]; then
        echo "This will stop the following instances:"
        list_instances
        read -p "Are you sure? (yes/no): " confirm
        if [ "$confirm" != "yes" ]; then
            echo "Cancelled."
            return
        fi
    fi
    
    for instance_id in $instance_ids; do
        echo "Stopping instance: $instance_id"
        oci compute instance action \
            --instance-id "$instance_id" \
            --action STOP \
            --wait-for-state STOPPED \
            || echo -e "${RED}Failed to stop $instance_id${NC}"
    done
    
    echo -e "${GREEN}Stop operation completed${NC}"
}

function start_instances() {
    local instance_ids=$(get_instance_ids)
    
    if [ -z "$instance_ids" ]; then
        echo -e "${YELLOW}No instances found with prefix: $PREFIX${NC}"
        return
    fi
    
    echo -e "${YELLOW}Starting instances with prefix: $PREFIX${NC}"
    
    for instance_id in $instance_ids; do
        echo "Starting instance: $instance_id"
        oci compute instance action \
            --instance-id "$instance_id" \
            --action START \
            --wait-for-state RUNNING \
            || echo -e "${RED}Failed to start $instance_id${NC}"
    done
    
    echo -e "${GREEN}Start operation completed${NC}"
}

function terminate_instances() {
    local instance_ids=$(get_instance_ids)
    
    if [ -z "$instance_ids" ]; then
        echo -e "${YELLOW}No instances found with prefix: $PREFIX${NC}"
        return
    fi
    
    echo -e "${RED}WARNING: This will permanently delete instances!${NC}"
    
    if [ "$AUTO_CONFIRM" != "true" ]; then
        echo "This will terminate the following instances:"
        list_instances
        read -p "Are you ABSOLUTELY sure? Type 'DELETE' to confirm: " confirm
        if [ "$confirm" != "DELETE" ]; then
            echo "Cancelled."
            return
        fi
    fi
    
    for instance_id in $instance_ids; do
        echo "Terminating instance: $instance_id"
        oci compute instance terminate \
            --instance-id "$instance_id" \
            --force \
            || echo -e "${RED}Failed to terminate $instance_id${NC}"
    done
    
    echo -e "${GREEN}Terminate operation completed${NC}"
}

function show_status() {
    echo -e "${YELLOW}Instance Status:${NC}"
    echo ""
    
    local instance_ids=$(get_instance_ids)
    
    if [ -z "$instance_ids" ]; then
        echo -e "${YELLOW}No instances found with prefix: $PREFIX${NC}"
        return
    fi
    
    for instance_id in $instance_ids; do
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        oci compute instance get --instance-id "$instance_id" \
            --query 'data.{Name:"display-name",State:"lifecycle-state",Shape:shape,OCPUS:"shape-config.ocpus",Memory:"shape-config"."memory-in-gbs",Created:"time-created"}' \
            || echo -e "${RED}Failed to get status for $instance_id${NC}"
        
        # Get public IP if available
        echo "Public IPs:"
        oci compute instance list-vnics --instance-id "$instance_id" \
            --query 'data[*]."public-ip"' \
            --raw-output | grep -v null || echo "  None"
        echo ""
    done
}

# Parse arguments
COMMAND=""
COMPARTMENT_ID=""
PREFIX="instance-"
AUTO_CONFIRM="false"

if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

COMMAND=$1
shift

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--compartment-id)
            COMPARTMENT_ID="$2"
            shift 2
            ;;
        -p|--prefix)
            PREFIX="$2"
            shift 2
            ;;
        -y|--yes)
            AUTO_CONFIRM="true"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Validate
check_oci_cli

if [ -z "$COMPARTMENT_ID" ]; then
    echo -e "${RED}Error: Compartment ID is required${NC}"
    print_usage
    exit 1
fi

# Execute command
case $COMMAND in
    list)
        list_instances
        ;;
    stop)
        stop_instances
        ;;
    start)
        start_instances
        ;;
    terminate)
        terminate_instances
        ;;
    status)
        show_status
        ;;
    *)
        echo -e "${RED}Unknown command: $COMMAND${NC}"
        print_usage
        exit 1
        ;;
esac
