package database

const bucket_user = "users"
const bucket_mailbox = "mailboxes"
const bucket_message = "messages"

const bucket_user_mailbox = "users.mailbox"
const bucket_mailbox_message = "mailboxes.message"

const index_user_username = "users.username"
const index_mailbox_user_id_name = "mailbox.user_id_name"
const index_message_mailbox_uid = "messages.mailbox_uid"

const counter_uidvalidity = "counters.uidvalidity"

const counter_uid = "counters.uid"
