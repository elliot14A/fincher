#!/usr/bin/env bash
set -euo pipefail

# Fincher GCP GCE VM Provisioning Script
# Creates an e2-medium instance with a 30GB Persistent Disk and installs Docker.
# Automatically tries alternate zones in India (Mumbai, Delhi) and Singapore if capacity is exhausted.

PROJECT_ID="${PROJECT_ID:-$(gcloud config get-value project)}"
INSTANCE_NAME="${INSTANCE_NAME:-fincher-server}"
MACHINE_TYPE="${MACHINE_TYPE:-e2-medium}"
DISK_SIZE="${DISK_SIZE:-30GB}"

CANDIDATE_ZONES=(
  "${ZONE:-asia-south1-b}"
  "asia-south1-c"
  "asia-south1-a"
  "asia-south2-a"
  "asia-south2-b"
  "asia-south2-c"
  "asia-southeast1-b"
  "asia-southeast1-a"
)

SELECTED_ZONE=""

for z in "${CANDIDATE_ZONES[@]}"; do
  echo "==> Trying to create '${INSTANCE_NAME}' in zone '${z}' (${MACHINE_TYPE}, ${DISK_SIZE} disk)..."
  if gcloud compute instances create "${INSTANCE_NAME}" \
    --project="${PROJECT_ID}" \
    --zone="${z}" \
    --machine-type="${MACHINE_TYPE}" \
    --network-interface="network-tier=PREMIUM,subnet=default" \
    --maintenance-policy="MIGRATE" \
    --scopes="https://www.googleapis.com/auth/cloud-platform" \
    --tags="http-server,https-server,fincher" \
    --create-disk="auto-delete=yes,boot=yes,image=projects/debian-cloud/global/images/family/debian-12,mode=rw,size=${DISK_SIZE},type=pd-balanced" \
    --quiet 2>/tmp/gce_err.log; then
    SELECTED_ZONE="${z}"
    echo "==> Successfully created instance in zone: ${SELECTED_ZONE}"
    break
  else
    echo "    Zone '${z}' unavailable (capacity exhausted), trying next zone..."
    cat /tmp/gce_err.log | grep -E "(ZONE_RESOURCE_POOL_EXHAUSTED|RESOURCE_AVAILABILITY|does not have enough resources)" || true
  fi
done

if [ -z "${SELECTED_ZONE}" ]; then
  echo "ERROR: Could not find available capacity in any tested zone. Check your GCP quota or try again in a few minutes."
  exit 1
fi

echo "==> Ensuring firewall rules for HTTP/HTTPS traffic..."
gcloud compute firewall-rules create allow-fincher-http \
  --project="${PROJECT_ID}" \
  --direction=INGRESS \
  --priority=1000 \
  --network=default \
  --action=ALLOW \
  --rules=tcp:80,tcp:443,tcp:8080 \
  --source-ranges=0.0.0.0/0 \
  --target-tags=fincher \
  --quiet || true

echo "==> Waiting 10s for VM SSH daemon to initialize..."
sleep 10

EXTERNAL_IP=$(gcloud compute instances describe "${INSTANCE_NAME}" --project="${PROJECT_ID}" --zone="${SELECTED_ZONE}" --format='get(networkInterfaces[0].accessConfigs[0].natIP)')

echo ""
echo "=========================================================="
echo " Fincher VM Provisioned for NixOS Bootstrap!"
echo " Instance:    ${INSTANCE_NAME}"
echo " Zone:        ${SELECTED_ZONE}"
echo " External IP: ${EXTERNAL_IP}"
echo ""
echo " Run nixos-anywhere to install:"
echo "   nix develop -c nixos-anywhere --flake .#fincher-prod -i ~/.ssh/google_compute_engine kiwi@${EXTERNAL_IP}"
echo "=========================================================="
