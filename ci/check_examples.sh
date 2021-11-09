echo "Checking examples..."

for example_directory in examples/*/; do
  echo "⌛ Checking $example_directory..."
  ./goal -c "$example_directory/goal.yaml"
  echo "✅ $example_directory ok"
  echo
done
