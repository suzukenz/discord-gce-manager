#!/bin/bash
if [ $# -eq 0 ] ; then
  # Use environments when no arguments, for Flexible Deployed Env.
  echo "No arguments supplied"
  "$HOME/$APP_NAME" -p "${PROJECT_ID}" -t "${DISCORD_TOKEN}" -d "${DISCORD_WEBHOOK}"
else
  "$HOME/$APP_NAME" ${@+"$@"}
fi
