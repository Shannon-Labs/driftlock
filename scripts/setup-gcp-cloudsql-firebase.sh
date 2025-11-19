#!/bin/bash

# Driftlock GCP Cloud SQL + Firebase Setup Script
# This script creates Cloud SQL database and Firebase Auth configuration

set -e

PROJECT_ID="driftlock"
REGION="us-central1"

echo "🚀 Setting up Driftlock with Cloud SQL + Firebase Auth..."
echo "Project: $PROJECT_ID"
echo "Region: $REGION"
echo ""

# Function to enable required APIs
enable_apis() {
    echo "🔧 Enabling required GCP APIs..."

    apis=(
        "sqladmin.googleapis.com"
        "secretmanager.googleapis.com"
        "run.googleapis.com"
        "cloudbuild.googleapis.com"
        "firebase.googleapis.com"
        "identitytoolkit.googleapis.com"
        "iam.googleapis.com"
    )

    for api in "${apis[@]}"; do
        echo "Enabling $api..."
        gcloud services enable "$api" --project="$PROJECT_ID"
    done

    echo "✅ APIs enabled"
}

# Function to set up Cloud SQL
setup_cloudsql() {
    echo "🗄️  Setting up Cloud SQL database..."

    # Check if instance already exists
    if gcloud sql instances describe driftlock-db --project="$PROJECT_ID" &>/dev/null; then
        echo "✅ Cloud SQL instance already exists"
    else
        echo "Creating Cloud SQL instance..."
        gcloud sql instances create driftlock-db \
            --project="$PROJECT_ID" \
            --database-version=POSTGRES_15 \
            --tier=db-custom-4-16384 \
            --region="$REGION" \
            --storage-size=100GB \
            --storage-type=SSD \
            --backup-start-time=02:00 \
            --availability-type=REGIONAL \
            --no-assign-ip

        echo "✅ Cloud SQL instance created"
    fi

    # Create database if it doesn't exist
    if ! gcloud sql databases describe driftlock --instance=driftlock-db --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating database..."
        gcloud sql databases create driftlock --instance=driftlock-db --project="$PROJECT_ID"
        echo "✅ Database created"
    else
        echo "✅ Database already exists"
    fi

    # Create user if it doesn't exist
    if ! gcloud sql users describe driftlock_user --instance=driftlock-db --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating database user..."

        # Generate secure password
        db_password=$(openssl rand -base64 32)

        gcloud sql users create driftlock_user \
            --instance=driftlock-db \
            --password="$db_password" \
            --project="$PROJECT_ID"

        # Save password to temporary file
        echo "$db_password" > /tmp/cloudsql-db-password.txt
        chmod 600 /tmp/cloudsql-db-password.txt

        echo "✅ Database user created"
        echo "💾 Password saved to /tmp/cloudsql-db-password.txt"
    else
        echo "✅ Database user already exists"
    fi
}

# Function to set up Firebase Auth
setup_firebase_auth() {
    echo "🔥 Setting up Firebase Authentication..."

    # Initialize Firebase if not already done
    if [ ! -f ".firebaserc" ]; then
        echo "Initializing Firebase..."
        firebase use "$PROJECT_ID"
        echo "✅ Firebase initialized"
    else
        echo "✅ Firebase already initialized"
    fi

    # Enable Firebase Auth
    echo "Enabling Firebase Authentication..."
    gcloud firebase projects enable auth --project="$PROJECT_ID"

    # Configure Firebase Auth settings
    echo "Configuring Firebase Auth settings..."

    # Enable Email/Password provider
    gcloud identitytoolkit config update \
        --project="$PROJECT_ID" \
        --sign-in-email=enabled \
        --anonymous=disabled

    echo "✅ Firebase Auth enabled"
    echo ""
    echo "📋 Next steps for Firebase Auth configuration:"
    echo "1. Go to https://console.firebase.google.com/project/$PROJECT_ID/authentication"
    echo "2. Configure sign-in providers (Email/Password, Google, etc.)"
    echo "3. Set up email templates for verification and password reset"
    echo "4. Configure authorized domains for your application"
}

# Function to get database connection string
get_db_connection_string() {
    echo "🔗 Generating database connection string..."

    # Get password from file or prompt user
    if [ -f "/tmp/cloudsql-db-password.txt" ]; then
        db_password=$(cat /tmp/cloudsql-db-password.txt)
    else
        echo "Enter your Cloud SQL database password:"
        read -s db_password
        echo ""
    fi

    # Generate connection string
    connection_string="postgresql://driftlock_user:$db_password@127.0.0.1:5432/driftlock?host=/cloudsql/$PROJECT_ID:$REGION:driftlock-db"

    echo "$connection_string"
    echo "$connection_string" > /tmp/cloudsql-connection-string.txt
    chmod 600 /tmp/cloudsql-connection-string.txt

    echo "💾 Connection string saved to /tmp/cloudsql-connection-string.txt"
}

# Function to generate Firebase service account key
generate_firebase_service_account() {
    echo "🔑 Generating Firebase service account key..."

    # Create service account for Firebase Admin
    sa_email="firebase-admin-sdk@$PROJECT_ID.iam.gserviceaccount.com"

    if ! gcloud iam service-accounts describe "$sa_email" --project="$PROJECT_ID" &>/dev/null; then
        echo "Creating Firebase Admin service account..."
        gcloud iam service-accounts create "$sa_email" \
            --display-name="Firebase Admin SDK Service Account" \
            --project="$PROJECT_ID"
    fi

    # Generate key file
    key_file="firebase-service-account-key.json"

    if [ ! -f "$key_file" ]; then
        gcloud iam service-accounts keys create "$key_file" \
            --iam-account="$sa_email" \
            --project="$PROJECT_ID"

        chmod 600 "$key_file"
        echo "✅ Firebase service account key created: $key_file"
    else
        echo "✅ Firebase service account key already exists"
    fi

    # Grant required permissions
    echo "Granting permissions to service account..."

    # Firebase Admin SDK roles
    roles=(
        "firebase.admin"
        "identitytoolkit.admin"
    )

    for role in "${roles[@]}"; do
        gcloud projects add-iam-policy-binding "$PROJECT_ID" \
            --member="serviceAccount:$sa_email" \
            --role="roles/$role" \
            --condition=None
    done

    echo "✅ Service account permissions granted"
}

# Function to create GCP secrets
create_secrets() {
    echo "🔐 Creating GCP secrets..."

    # Get connection string
    connection_string=$(get_db_connection_string)

    # Create secrets one by one

    echo "Creating driftlock-db-url secret..."
    echo -n "$connection_string" | gcloud secrets create driftlock-db-url \
        --data-file=- \
        --project="$PROJECT_ID" \
        --description="Cloud SQL PostgreSQL connection string" \
        --replication-policy="automatic"

    echo "Creating driftlock-license-key secret..."
    echo -n "dev-mode" | gcloud secrets create driftlock-license-key \
        --data-file=- \
        --project="$PROJECT_ID" \
        --description="License key for Driftlock (dev-mode for testing)" \
        --replication-policy="automatic"

    echo "Creating firebase-service-account-key secret..."
    if [ -f "firebase-service-account-key.json" ]; then
        gcloud secrets create firebase-service-account-key \
            --data-file="firebase-service-account-key.json" \
            --project="$PROJECT_ID" \
            --description="Firebase Admin SDK service account key" \
            --replication-policy="automatic"
    fi

    echo "Creating driftlock-api-key secret..."
    api_key="sk-drltlck-$(openssl rand -hex 16)"
    echo -n "$api_key" | gcloud secrets create driftlock-api-key \
        --data-file=- \
        --project="$PROJECT_ID" \
        --description="Admin API key for service management" \
        --replication-policy="automatic"

    echo "✅ Core secrets created"

    # Prompt for other secrets
    echo ""
    echo "📧 Configure SendGrid (optional):"
    read -p "Enter SendGrid API key (or press Enter to skip): " sendgrid_key
    if [ -n "$sendgrid_key" ]; then
        echo -n "$sendgrid_key" | gcloud secrets create sendgrid-api-key \
            --data-file=- \
            --project="$PROJECT_ID" \
            --description="SendGrid API key for email delivery" \
            --replication-policy="automatic"
    fi

    echo ""
    echo "💳 Configure Stripe:"
    read -p "Enter Stripe secret key (sk_test_...): " stripe_key
    if [ -n "$stripe_key" ]; then
        echo -n "$stripe_key" | gcloud secrets create stripe-secret-key \
            --data-file=- \
            --project="$PROJECT_ID" \
            --description="Stripe secret key for payments" \
            --replication-policy="automatic"
    fi

    read -p "Enter Stripe price ID for Pro plan: " stripe_price_id
    if [ -n "$stripe_price_id" ]; then
        echo -n "$stripe_price_id" | gcloud secrets create stripe-price-id-pro \
            --data-file=- \
            --project="$PROJECT_ID" \
            --description="Stripe price ID for Pro plan" \
            --replication-policy="automatic"
    fi
}

# Function to run database migrations
run_migrations() {
    echo "🔄 Running database migrations..."

    # Check if migration files exist
    if [ -d "api/migrations" ]; then
        echo "Found migration files in api/migrations/"

        # Connect to Cloud SQL and run migrations
        # Using gcloud sql connect to run SQL files
        for migration_file in api/migrations/*.sql; do
            if [ -f "$migration_file" ]; then
                echo "Running migration: $migration_file"
                gcloud sql connect driftlock-db --user=driftlock_user --project="$PROJECT_ID" < "$migration_file"
            fi
        done

        echo "✅ Database migrations completed"
    else
        echo "⚠️  No migration files found at api/migrations/"
        echo "You may need to run migrations manually via:"
        echo "1. Google Cloud Console > Cloud SQL"
        echo "2. Select driftlock-db instance"
        echo "3. Go to 'Databases' > 'driftlock' > 'Query'"
    fi
}

# Function to configure Cloud Run for Cloud SQL
configure_cloudrun() {
    echo "☁️  Configuring Cloud Run for Cloud SQL access..."

    # Get the Cloud SQL connection instance name
    instance_connection_name="$PROJECT_ID:$REGION:driftlock-db"

    echo "Cloud SQL instance connection name: $instance_connection_name"
    echo ""
    echo "This will be automatically configured in your Cloud Run deployment via:"
    echo "1. --add-cloudsql-instances flag in gcloud run deploy"
    echo "2. DATABASE_URL secret in Cloud Run environment"
    echo ""
    echo "The deployment script will handle this automatically."
}

# Function to provide next steps
next_steps() {
    echo ""
    echo "🎯 Cloud SQL + Firebase Auth Setup Complete!"
    echo "=========================================="
    echo ""
    echo "📋 Configuration Summary:"
    echo "- Cloud SQL Instance: driftlock-db"
    echo "- Database: driftlock"
    echo "- User: driftlock_user"
    echo "- Firebase Auth: Enabled"
    echo "- Service Account: Created and configured"
    echo ""
    echo "🔑 Important Files:"
    echo "- Firebase Service Account: firebase-service-account-key.json"
    echo "- DB Password: /tmp/cloudsql-db-password.txt"
    echo "- DB Connection String: /tmp/cloudsql-connection-string.txt"
    echo ""
    echo "🌐 Console Links:"
    echo "- Cloud SQL: https://console.cloud.google.com/sql/instances/driftdb-db"
    echo "- Firebase Auth: https://console.firebase.google.com/project/$PROJECT_ID/authentication"
    echo "- Secret Manager: https://console.cloud.google.com/security/secret-manager"
    echo ""
    echo "🚀 Next Steps:"
    echo "1. Configure Firebase Auth providers in the console"
    echo "2. Set up custom email templates in Firebase"
    echo "3. Add Stripe and SendGrid secrets when ready"
    echo "4. Run: ./scripts/deploy-production.sh"
    echo ""
    echo "🧪 Test the setup:"
    echo "./scripts/test-deployment-complete.sh"
}

# Main execution
echo "Setting up Driftlock with Cloud SQL + Firebase Auth..."
echo ""

enable_apis
setup_cloudsql
setup_firebase_auth
generate_firebase_service_account
create_secrets
run_migrations
configure_cloudrun
next_steps