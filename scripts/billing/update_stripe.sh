#!/bin/bash
set -e

# Archive old prices
echo "Archiving old prices..."
prices=(
  "price_1SW89dL4rhSbUSqAZ42gNxZb"
  "price_1SO7RbL4rhSbUSqAYE8UPlmL"
  "price_1SNmfBL4rhSbUSqAS5M5ai9N"
  "price_1SNmf8L4rhSbUSqAQxgW1yxE"
  "price_1SMi3xL4rhSbUSqA7TFubIed"
  "price_1SMi3oL4rhSbUSqAO3hkhM0Y"
  "price_1SMhshL4rhSbUSqAyHfhWUSQ"
  "price_1SMhsZL4rhSbUSqA51lWvPlQ"
  "price_1SMhhuL4rhSbUSqAM1l1hEEC"
  "price_1SMhhUL4rhSbUSqAEKMoro2d"
)

for price in "${prices[@]}"; do
  echo "Archiving $price..."
  stripe prices update "$price" --active=false
done

echo "Creating new products and prices..."

# Create Radar Plan ($15/month)
echo "Creating Driftlock Radar..."
stripe products create --name "Driftlock Radar" --description "Entry level AI plan" > radar_prod.json 2>radar_prod.err
if [ ! -s radar_prod.json ]; then
  echo "Error creating product:"
  cat radar_prod.err
  exit 1
fi
radar_prod=$(jq -r .id radar_prod.json)
echo "Product ID: $radar_prod"

stripe prices create --product "$radar_prod" --unit-amount 1500 --currency usd --recurring.interval month > radar_price.json 2>radar_price.err
radar_price=$(jq -r .id radar_price.json)
echo "Radar Price ID: $radar_price"

# Create Tensor Plan ($100/month) - Renamed from Lock
echo "Creating Driftlock Tensor..."
stripe products create --name "Driftlock Tensor" --description "Pro level AI plan with Claude Sonnet 4.5" > tensor_prod.json
tensor_prod=$(jq -r .id tensor_prod.json)
stripe prices create --product "$tensor_prod" --unit-amount 10000 --currency usd --recurring.interval month > tensor_price.json
tensor_price=$(jq -r .id tensor_price.json)
echo "Tensor Price ID: $tensor_price"

# Create Orbit Plan ($499/month)
echo "Creating Driftlock Orbit..."
stripe products create --name "Driftlock Orbit" --description "Enterprise level AI plan" > orbit_prod.json
orbit_prod=$(jq -r .id orbit_prod.json)
stripe prices create --product "$orbit_prod" --unit-amount 49900 --currency usd --recurring.interval month > orbit_price.json
orbit_price=$(jq -r .id orbit_price.json)
echo "Orbit Price ID: $orbit_price"

# Update .env file
echo "Updating .env file..."
# We need to replace or append the price IDs
# Assuming .env exists and we can use sed
if [ -f .env ]; then
  sed -i '' "s/STRIPE_PRICE_ID_BASIC=.*/STRIPE_PRICE_ID_BASIC=$radar_price/" .env
  sed -i '' "s/STRIPE_PRICE_ID_PRO=.*/STRIPE_PRICE_ID_PRO=$tensor_price/" .env
  # Add Enterprise if not exists
  if grep -q "STRIPE_PRICE_ID_ENTERPRISE" .env; then
    sed -i '' "s/STRIPE_PRICE_ID_ENTERPRISE=.*/STRIPE_PRICE_ID_ENTERPRISE=$orbit_price/" .env
  else
    echo "STRIPE_PRICE_ID_ENTERPRISE=$orbit_price" >> .env
  fi
else
  echo "Warning: .env file not found. Please update it manually."
fi

# Display summary
echo ""
echo "=== Stripe Products Updated ==="
echo "Radar (Basic): $radar_price"
echo "Tensor (Pro):  $tensor_price"
echo "Orbit (Enterprise): $orbit_price"
echo ""
echo "Remember to update Google Secret Manager:"
echo "  gcloud secrets versions add stripe-price-id-basic --data-file=/dev/stdin <<< \"$radar_price\""
echo "  gcloud secrets versions add stripe-price-id-pro --data-file=/dev/stdin <<< \"$tensor_price\""
echo "  gcloud secrets versions add stripe-price-id-enterprise --data-file=/dev/stdin <<< \"$orbit_price\""

echo "Done!"
