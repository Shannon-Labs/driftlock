# Cloudflare R2 Storage Configuration for Driftlock
# This file defines the infrastructure configuration for Cloudflare R2 storage

resource "cloudflare_r2_bucket" "driftlock_uploads" {
  account_id = var.cloudflare_account_id
  name       = "driftlock-file-uploads"
  location   = "auto"  # Let Cloudflare optimize for performance
}

# CORS Configuration for the R2 bucket
resource "cloudflare_r2_bucket_cors_configuration" "driftlock_cors" {
  account_id = var.cloudflare_account_id
  bucket_name = cloudflare_r2_bucket.driftlock_uploads.name

  cors_rules {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allowed_origins = [
      "https://driftlock.net",
      "https://www.driftlock.net",
      "https://api.driftlock.net"
    ]
    expose_headers = ["ETag", "Content-Length"]
    max_age_seconds = 3600
  }
}

# R2 Bucket Public Access Settings (if needed)
resource "cloudflare_r2_bucket_public_access_block" "driftlock_public_access" {
  account_id = var.cloudflare_account_id
  bucket_name = cloudflare_r2_bucket.driftlock_uploads.name

  # Block all public access by default
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Variables needed
variable "cloudflare_account_id" {
  description = "Cloudflare Account ID"
  type        = string
}

# Outputs
output "r2_bucket_name" {
  description = "Name of the R2 bucket"
  value       = cloudflare_r2_bucket.driftlock_uploads.name
}

output "r2_bucket_endpoint" {
  description = "R2 bucket endpoint URL"
  value       = "${cloudflare_r2_bucket.driftlock_uploads.name}.r2.cloudflarestorage.com"
}