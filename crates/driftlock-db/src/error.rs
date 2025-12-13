//! Database error types

use thiserror::Error;

#[derive(Error, Debug)]
pub enum DbError {
    #[error("Record not found")]
    NotFound,

    #[error("Duplicate key: {0}")]
    DuplicateKey(String),

    #[error("Foreign key violation: {0}")]
    ForeignKeyViolation(String),

    #[error("Invalid verification token")]
    InvalidToken,

    #[error("Token expired")]
    TokenExpired,

    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("Database error: {0}")]
    Sqlx(sqlx::Error),

    #[error("Pool error: {0}")]
    Pool(String),

    #[error("Crypto error: {0}")]
    Crypto(String),
}

impl From<sqlx::Error> for DbError {
    fn from(err: sqlx::Error) -> Self {
        match &err {
            sqlx::Error::RowNotFound => DbError::NotFound,
            sqlx::Error::Database(db_err) => {
                if let Some(code) = db_err.code() {
                    match code.as_ref() {
                        "23505" => DbError::DuplicateKey(db_err.message().to_string()),
                        "23503" => DbError::ForeignKeyViolation(db_err.message().to_string()),
                        _ => DbError::Sqlx(err),
                    }
                } else {
                    DbError::Sqlx(err)
                }
            }
            _ => DbError::Sqlx(err),
        }
    }
}
