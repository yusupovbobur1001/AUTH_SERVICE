CURRENT_DIR=$(pwd)
OUTPUT_DIR="${CURRENT_DIR}"

mkdir -p "$OUTPUT_DIR" || { echo "Failed to create output directory"; exit 1; }
 

for dir in "${CURRENT_DIR}/PRODUCTS"/*; do
    # Skip non-directories 
    [ -d "$dir" ] || continue

    for proto_file in "$dir"/*.proto; do
        echo "Compiling $proto_file..."
        protoc -I="${dir}" -I="${CURRENT_DIR}/PRODUCTS" -I /usr/local/include --go_out="$OUTPUT_DIR" --go-grpc_out="$OUTPUT_DIR" "$proto_file" || { echo "Failed to compile $proto_file"; exit 1; }
    done
done

echo "Compilation completed successfully."