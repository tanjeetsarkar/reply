# Remove ../../api/static directory if already present
if [ -d ../../api/static ]; then
    rm -rf ../../api/static
fi

# Check if ../../api/static directory exists
if [ ! -d ../../api/static ]; then
    # Create the static directory
    mkdir -p ../../api/static
fi

# Move contents of dist folder to ../../api/static
mv dist/* ../../api/static
