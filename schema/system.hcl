table "access_tokens" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("generate_ulid()")
  }
  column "name" {
    null = false
    type = text
  }
  column "person_id" {
    null = false
    type = uuid
  }
  column "value" {
    null = false
    type = text
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "expires_at" {
    null = true # We can have never expiring tokens
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "access_tokens_person_name_unique_key" {
    unique  = true
    columns = [column.person_id, column.name]
  }
  foreign_key "access_tokens_person_fkey" {
    columns     = [column.person_id]
    ref_columns = [table.people.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

table "event_queue" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("generate_ulid()")
  }
  column "name" {
    null = false
    type = text
  }
  column "properties" {
    null = true
    type = jsonb
  }
  column "error" {
    null = true
    type = text
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "last_attempt" {
    null = true
    type = timestamptz
  }
  column "attempts" {
    null    = true
    type    = integer
    default = 0
  }
  column "priority" {
    null    = false
    type    = integer
    default = 100
  }
  primary_key {
    columns = [column.id]
  }
  index "event_queue_name_properties" {
    unique  = true
    columns = [column.name, column.properties]
  }
  index "event_queue_properties" {
    type    = GIN
    columns = [column.properties]
  }
  index "event_queue_pop" {
    columns = [column.name, column.priority, column.created_at]
  }
}

table "integrations" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("generate_ulid()")
  }
  column "name" {
    null = false
    type = text
  }
  column "icon" {
    null = true
    type = text
  }
  column "spec" {
    null = true
    type = jsonb
  }
  column "source" {
    null = true
    type = text
  }
  column "created_by" {
    null = false
    type = uuid
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "integrations_created_by_fkey" {
    columns     = [column.created_by]
    ref_columns = [table.people.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}

table "job_history" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("generate_ulid()")
  }
  column "name" {
    null = true
    type = text
  }
  column "success_count" {
    null = true
    type = integer
  }
  column "error_count" {
    null = true
    type = integer
  }
  column "details" {
    null = true
    type = jsonb
  }
  column "hostname" {
    null = true
    type = text
  }
  column "duration_millis" {
    null = true
    type = integer
  }
  column "resource_type" {
    null = true
    type = text
  }
  column "resource_id" {
    null = true
    type = text
  }
  column "status" {
    null = true
    type = text
  }
  column "time_start" {
    null = true
    type = timestamptz
  }
  column "time_end" {
    null = true
    type = timestamptz
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
}
