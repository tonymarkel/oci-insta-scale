#!/bin/bash

# Define your compartment OCID and region
COMPARTMENT_OCID="ocid1.compartment.oc1..aaaaaaaaqrh2fgi5jycaurey53m64lhrou4peu6ffq7mgu2lua74iylb3ata"
REGION="us-ashburn-1"

echo "Listing unattached boot volumes in compartment: $COMPARTMENT_OCID in region: $REGION"

# List all boot volumes and filter for unattached ones using jq
# An unattached volume has a 'lifecycle-state' of 'AVAILABLE' and a null 'attachment-id' in the list output
UNATTACHED_BOOT_VOLUMES=$(oci bv boot-volume list --compartment-id "$COMPARTMENT_OCID" --region "$REGION" --all | jq -r '.data[] | select(."attachment-id" | not) | ."id"')

# Check if any unattached volumes were found
if [ -z "$UNATTACHED_BOOT_VOLUMES" ]; then
    echo "No unattached boot volumes found."
else
    echo "Found unattached boot volumes. Proceeding with deletion."
    # Loop through the list of unattached boot volume OCIDs and delete each one
    for BV_ID in $UNATTACHED_BOOT_VOLUMES; do
        echo "Deleting boot volume ID: $BV_ID"
        # Use the --force option to delete without a confirmation prompt
        oci bv boot-volume delete --boot-volume-id "$BV_ID" --force
        echo "Boot volume $BV_ID deleted."
    done
fi

echo "Deletion process complete."
