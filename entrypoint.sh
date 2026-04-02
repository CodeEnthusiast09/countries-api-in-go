#!/bin/sh
# The line above tells the OS to run this with /bin/sh (the shell).
# We use /bin/sh (not /bin/bash) because Alpine Linux only ships with sh.

# set -e means: "exit immediately if ANY command in this script fails".
# Without this, the script would continue to the next line even after an error.
# With this, if atlas migrate apply fails, the container exits — Railway sees
# the failure and reports it, rather than your broken app starting silently.
set -e

# =============================================================================
# STEP 1: Wait for the database to be ready
# =============================================================================
# Railway starts your app and its MySQL service roughly at the same time.
# Your app container might boot faster than MySQL is ready to accept connections.
# We use a retry loop to wait for MySQL to be reachable.
#
# How it works:
# - We try `atlas migrate status` which just checks the DB connection.
# - If it fails (DB not ready), we wait 3 seconds and try again.
# - If it keeps failing for more than MAX_RETRIES attempts, we give up and exit.
# - The "2>/dev/null" part suppresses error output during the waiting phase
#   so your logs aren't flooded with "connection refused" messages.

MAX_RETRIES=30
count=0

echo "Waiting for database to be ready..."

until atlas migrate status \
  --url "$ATLAS_DATABASE_URL" \
  --dir "file:///app/migrations" > /dev/null 2>&1; do

  count=$((count + 1))

  if [ "$count" -ge "$MAX_RETRIES" ]; then
    echo "ERROR: Database did not become ready after $MAX_RETRIES attempts."
    echo "Check that ATLAS_DATABASE_URL is set correctly in Railway."
    exit 1
  fi

  echo "Database not ready yet (attempt $count/$MAX_RETRIES). Retrying in 3s..."
  sleep 3
done

echo "Database is ready."

# =============================================================================
# STEP 2: Run pending migrations
# =============================================================================
# This runs `atlas migrate apply` — the same command you run locally, but
# pointing at the production database URL from your environment variables.
#
# Atlas is smart: it checks the atlas_schema_revisions table to know which
# migrations have already been applied, and only runs the new ones.
# On first deploy, it runs all migrations. On subsequent deploys, it only
# runs whatever is new.

echo "Applying database migrations..."

atlas migrate apply \
  --url "$ATLAS_DATABASE_URL" \
  --dir "file:///app/migrations"

echo "Migrations applied successfully."

# =============================================================================
# STEP 3: Start the application
# =============================================================================
# `exec` replaces the shell process with the server process.
# This is important: without exec, the shell stays running as process ID 1
# and the server runs as a child. With exec, the server becomes PID 1 directly,
# which means Railway's signals (like SIGTERM for graceful shutdown) go directly
# to your app rather than to the shell.

echo "Starting server..."
exec /app/server
