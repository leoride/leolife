CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS alert_field;
DROP TABLE IF EXISTS alert;
DROP TABLE IF EXISTS alert_type_field;
DROP TABLE IF EXISTS alert_type;

CREATE TABLE alert_type (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" text UNIQUE NOT NULL,
  description text NOT NULL
);
CREATE TABLE alert_type_field (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  alert_type_id uuid REFERENCES alert_type(id) NOT NULL,
  "name" text NOT NULL,
  "label" text NOT NULL,
  "type" text NOT NULL,
  "default" text,
  mandatory boolean NOT NULL DEFAULT false
);

CREATE TABLE alert (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name" text UNIQUE NOT NULL,
  description text NOT NULL,
  alert_type_id uuid REFERENCES alert_type(id) NOT NULL,
  start_timestamp timestamp with time zone NOT NULL,
  end_timestamp timestamp with time zone NOT NULL
);
CREATE TABLE alert_field (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  alert_id uuid REFERENCES alert(id) NOT NULL,
  alert_type_field_id uuid REFERENCES alert_type_field(id) NOT NULL,
  value text
);

INSERT INTO alert_type VALUES (
  uuid_generate_v4(),
  'birthday',
  'Birthday alert'
);
INSERT INTO alert_type_field VALUES (
  uuid_generate_v4(),
  (select id from alert_type where name = 'birthday'),
  'date',
  'Birthday date',
  'date',
  null,
  TRUE
);
INSERT INTO alert_type_field VALUES (
  uuid_generate_v4(),
  (select id from alert_type where name = 'birthday'),
  'name',
  'Name',
  'string',
  null,
  TRUE
);

INSERT INTO alert VALUES (
  uuid_generate_v4(),
  'Tom Birthday',
  'Tom Birthday is today!',
  (select id from alert_type where name = 'birthday'),
  '2015-08-15T00:00:00.000+00',
  '2015-08-16T00:00:00.000+00'
);

INSERT INTO alert_field VALUES (
  uuid_generate_v4(),
  (select id from alert where name = 'Tom Birthday'),
  (select id from alert_type_field where name = 'date'),
  '08/15/1989'
);

INSERT INTO alert_field VALUES (
  uuid_generate_v4(),
  (select id from alert where name = 'Tom Birthday'),
  (select id from alert_type_field where name = 'name'),
  'Tom Granger'
);