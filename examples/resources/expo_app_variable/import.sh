# Environment variables are scoped to an app, so the import id is
# <app_id>/<variable_id>. Both ids come from `eas project:info` and
# `eas env:list --environment <env> --format long`.
terraform import expo_app_variable.this 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111
