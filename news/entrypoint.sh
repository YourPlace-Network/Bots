#!/bin/sh
if [ -n "$CONFIG_JSON" ]; then
    printf '%s' "$CONFIG_JSON" > /app/config.json
fi
if [ -n "$WALLET_JSON" ]; then
    printf '%s' "$WALLET_JSON" > /app/wallet.json
fi
exec ./newsbot -config=/app/config.json -data=/app/data "$@"
