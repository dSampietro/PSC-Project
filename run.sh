#!/bin/sh

# Default values
FILE=""
MAX_DEPTH=1
SEQ=false
PRINT_SENTENCES=false

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    -file)
      FILE="$2"
      shift 2
      ;;
    -max_depth)
      MAX_DEPTH="$2"
      shift 2
      ;;
    -seq)
      SEQ=true
      shift 1
      ;;
    -print_sentences)
      PRINT_SENTENCES=true
      shift 1
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 -file <filename> -max_depth <num> [-seq] [-print_sentences]"
      exit 1
      ;;
  esac
done

# Run the Go program
MAIN="sketch2.go" 
CMD="go run $MAIN graph.go file_operation.go strategy.go 
    -file \"$FILE\" -max_depth $MAX_DEPTH"

if [ "$PRINT_SENTENCES" = true ]; then
  CMD="$CMD -print_sentences"
fi

if [ "$SEQ" = true ]; then
  CMD="$CMD -seq"
fi

eval $CMD
