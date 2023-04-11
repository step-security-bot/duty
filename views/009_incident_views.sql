DROP FUNCTION if exists lookup_component_incidents;

CREATE OR REPLACE FUNCTION lookup_component_incidents(component_id text)
RETURNS TABLE (
    id UUID
) AS $$
BEGIN
    RETURN QUERY
        SELECT incidents.id FROM incidents WHERE incidents.id IN (
            SELECT incident_id FROM hypotheses WHERE hypotheses.id IN (
                SELECT hypothesis_id FROM evidences WHERE evidences.component_id = $1::UUID
            )
        );
END;
$$
language plpgsql;


-- incidents_by_component
DROP VIEW IF EXISTS incidents_by_component;
CREATE OR REPLACE VIEW incidents_by_component AS
  SELECT DISTINCT incidents.title,incidents.id, incidents.created_at, incidents."type", incidents.status, incidents.severity, component_id FROM evidences
  INNER join hypotheses on evidences.hypothesis_id = hypotheses.id
  INNER JOIN incidents on hypotheses.incident_id = incidents.id
  WHERE component_id is not null;


--incidents_by_config
DROP VIEW IF EXISTS incidents_by_config;
CREATE OR REPLACE VIEW incidents_by_config AS
  SELECT DISTINCT incidents.title, incidents.id, incidents.created_at, incidents."type", incidents.status,  incidents.severity, config_id FROM evidences
  INNER join hypotheses on evidences.hypothesis_id = hypotheses.id
  INNER JOIN incidents on hypotheses.incident_id = incidents.id
  WHERE evidences.config_id is not null;

-- incident_summary_by_component
DROP VIEW IF EXISTS incident_summary_by_component;
CREATE OR REPLACE VIEW incident_summary_by_component AS
  WITH type_summary AS (
      SELECT summary.id, summary.type, json_object_agg(f.k, f.v) as json
      FROM (
          SELECT evidences.component_id AS id, incidents.type, json_build_object(severity, count(*)) AS severity_agg
          FROM incidents
          INNER JOIN hypotheses ON hypotheses.incident_id = incidents.id
          INNER JOIN evidences ON evidences.hypothesis_id = hypotheses.id
          WHERE (incidents.resolved IS NULL AND incidents.closed IS NULL and evidences.component_id IS NOT NULL
      )
      GROUP BY incidents.severity, incidents.type, evidences.component_id)
      AS summary, json_each(summary.severity_agg) AS f(k,v) GROUP BY summary.type, summary.id
  )

  SELECT id, jsonb_object_agg(key, value) as incidents FROM (select id, json_object_agg(type,json) incidents from type_summary group by id, type) i, json_each(incidents) group by id;

-- incident_summary VIEW
DROP VIEW IF EXISTS incident_summary;
CREATE OR REPLACE VIEW incident_summary AS 
  SELECT
    incidents.id,
    incidents.title,
    incidents.severity,
    incidents.type,
    incidents.status,
    incidents.created_at,
    incidents.updated_at,
    jsonb_build_object(
      'id', people.id,
      'avatar', people.avatar,
      'name', people.name
    ) AS commander,
    jsonb_agg(
      DISTINCT jsonb_build_object(
        'id', responder_person.id,
        'avatar', responder_person.avatar,
        'name', responder_person.name
      )
    ) FILTER (WHERE responder_person.id IS NOT NULL) AS distinct_responders,
    jsonb_agg(
      DISTINCT jsonb_build_object(
        'id', commenter.id,
        'avatar', commenter.avatar,
        'name', commenter.name
      )
    ) FILTER (WHERE commenter.id IS NOT NULL) AS distinct_commenters
  FROM
    incidents
    LEFT JOIN people ON incidents.commander_id = people.id
    LEFT JOIN responders ON incidents.id = responders.incident_id LEFT JOIN people responder_person ON responders.person_id = responder_person.id
    LEFT JOIN comments ON incidents.id = comments.incident_id LEFT JOIN people commenter ON comments.created_by = commenter.id
  GROUP BY
    incidents.id,
    people.id;