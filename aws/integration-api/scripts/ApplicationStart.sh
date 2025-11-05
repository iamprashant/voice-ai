#!/bin/bash
export ENV_PATH="/opt/app/backend-app/artifacts/integration-api/env"
/opt/app/backend-app/artifacts/integration-api/integration-api.0.0.1 > /dev/null 2> /dev/null < /dev/null &
