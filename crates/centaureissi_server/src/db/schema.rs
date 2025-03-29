// @generated automatically by Diesel CLI.

diesel::table! {
    messages (id) {
        id -> Integer,
        user_id -> Integer,
        content_hash -> Text,
        created_at -> Nullable<Timestamp>,
        updated_at -> Nullable<Timestamp>,
    }
}

diesel::table! {
    user_tokens (id) {
        id -> Integer,
        user_id -> Integer,
        token -> Text,
        revoked_at -> Nullable<Timestamp>,
        created_at -> Nullable<Timestamp>,
        updated_at -> Nullable<Timestamp>,
    }
}

diesel::table! {
    users (id) {
        id -> Integer,
        username -> Text,
        password -> Text,
        created_at -> Nullable<Timestamp>,
        updated_at -> Nullable<Timestamp>,
    }
}

diesel::allow_tables_to_appear_in_same_query!(
    messages,
    user_tokens,
    users,
);
