#!/bin/bash
set -e

echo "=== Migration Job Testing Script ==="
echo

# Variables
NAMESPACE="saas-bff"
JOB_NAME="prod-affiliate-backend-migration"

echo "Testing migration job in namespace: $NAMESPACE"
echo

# Function to check if job exists
check_job_exists() {
    kubectl get job $JOB_NAME -n $NAMESPACE >/dev/null 2>&1
}

# Function to delete existing job
delete_existing_job() {
    if check_job_exists; then
        echo "Deleting existing migration job..."
        kubectl delete job $JOB_NAME -n $NAMESPACE
        echo "✅ Existing job deleted"
    fi
}

# Function to create and monitor job
run_migration_job() {
    echo "Creating migration job..."
    kubectl apply -f /workspace/k8s/overlays/prod/migration-job.yaml
    echo "✅ Migration job created"
    
    echo "Waiting for job to start..."
    sleep 5
    
    # Get pod name
    POD_NAME=$(kubectl get pods -n $NAMESPACE -l job-name=$JOB_NAME -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
    
    if [ -z "$POD_NAME" ]; then
        echo "❌ No pod found for migration job"
        return 1
    fi
    
    echo "Migration pod: $POD_NAME"
    echo
    
    # Monitor job status
    echo "Monitoring job progress..."
    echo "Use Ctrl+C to stop monitoring (job will continue running)"
    echo
    
    # Follow logs
    kubectl logs -f $POD_NAME -n $NAMESPACE -c migration || true
    
    echo
    echo "Checking final job status..."
    kubectl get job $JOB_NAME -n $NAMESPACE
    
    # Check if job succeeded
    JOB_STATUS=$(kubectl get job $JOB_NAME -n $NAMESPACE -o jsonpath='{.status.conditions[?(@.type=="Complete")].status}' 2>/dev/null || echo "")
    
    if [ "$JOB_STATUS" = "True" ]; then
        echo "✅ Migration job completed successfully!"
        return 0
    else
        echo "❌ Migration job failed or is still running"
        echo "Check logs with: kubectl logs $POD_NAME -n $NAMESPACE -c migration"
        echo "Check job status with: kubectl describe job $JOB_NAME -n $NAMESPACE"
        return 1
    fi
}

# Function to show troubleshooting info
show_troubleshooting() {
    echo
    echo "=== Troubleshooting Information ==="
    echo
    
    echo "1. Check secrets exist:"
    kubectl get secrets -n $NAMESPACE | grep -E "(saas-bff-db-credentials|saas-bff-secrets)" || echo "❌ Required secrets not found"
    
    echo
    echo "2. Check service account:"
    kubectl get serviceaccount default -n $NAMESPACE -o yaml | grep -A 5 annotations || echo "❌ No annotations found"
    
    echo
    echo "3. Recent events:"
    kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp' | tail -10
    
    if check_job_exists; then
        POD_NAME=$(kubectl get pods -n $NAMESPACE -l job-name=$JOB_NAME -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
        if [ -n "$POD_NAME" ]; then
            echo
            echo "4. Pod description:"
            kubectl describe pod $POD_NAME -n $NAMESPACE
        fi
    fi
}

# Main execution
case "${1:-run}" in
    "run")
        delete_existing_job
        if run_migration_job; then
            echo "✅ Migration test completed successfully"
        else
            echo "❌ Migration test failed"
            show_troubleshooting
            exit 1
        fi
        ;;
    "clean")
        delete_existing_job
        echo "✅ Migration job cleaned up"
        ;;
    "status")
        if check_job_exists; then
            kubectl get job $JOB_NAME -n $NAMESPACE
            kubectl describe job $JOB_NAME -n $NAMESPACE
        else
            echo "No migration job found"
        fi
        ;;
    "logs")
        if check_job_exists; then
            POD_NAME=$(kubectl get pods -n $NAMESPACE -l job-name=$JOB_NAME -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
            if [ -n "$POD_NAME" ]; then
                kubectl logs $POD_NAME -n $NAMESPACE -c migration
            else
                echo "No pod found for migration job"
            fi
        else
            echo "No migration job found"
        fi
        ;;
    "troubleshoot")
        show_troubleshooting
        ;;
    *)
        echo "Usage: $0 [run|clean|status|logs|troubleshoot]"
        echo "  run         - Delete existing job and run new migration (default)"
        echo "  clean       - Delete existing migration job"
        echo "  status      - Show migration job status"
        echo "  logs        - Show migration job logs"
        echo "  troubleshoot - Show troubleshooting information"
        exit 1
        ;;
esac