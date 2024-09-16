CREATE TABLE IF NOT EXISTS employee (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      username VARCHAR(50) UNIQUE NOT NULL,
      first_name VARCHAR(50),
      last_name VARCHAR(50),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE  organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE IF NOT EXISTS organization (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      name VARCHAR(100) NOT NULL,
      description TEXT,
      type organization_type,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
      id SERIAL PRIMARY KEY,
      organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
      user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TABLE IF NOT EXISTS tender (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(100) NOT NULL,
            description TEXT,
            type service_type,
            status status DEFAULT 'Created',
            organization_id UUID DEFAULT gen_random_uuid(),
            version INT DEFAULT 1 CHECK ( version >= 1 ),
            created_by VARCHAR(50) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION get_user(input VARCHAR(50)) RETURNS INT AS
$$
DECLARE user_id INT;
BEGIN
    SELECT id INTO user_id FROM employee WHERE username=input;
    IF user_id IS NULL THEN
            RAISE 'unknown user' USING ERRCODE='NO_USER_FOUND';
    END IF;
    RETURN user_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_tender(usr VARCHAR(50), tender_id UUID) RETURNS UUID AS
$$
DECLARE
uid UUID = get_user(usr);
BEGIN
    IF NOT EXISTS(SELECT 1 FROM tender WHERE id=tender_id) THEN
        RAISE 'tender does not exist' USING ERRCODE='NO_TENDER_FOUND';
    ELSEIF NOT EXISTS(SELECT 1 FROM tender JOIN organization_responsible o ON o.organization_id = tender.organization_id
                                           JOIN employee e on e.id = o.user_id WHERE user_id=uid) THEN
        RAISE 'no rights to tender' USING ERRCODE='ACCESS_DENIED';
END IF;
RETURN tender_id;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION update_tender(usr VARCHAR(50), tender_id UUID, nm VARCHAR(50), des TEXT, tp SERVICE_TYPE, st STATUS)  RETURNS tender AS
$$
DECLARE
    t tender;
BEGIN
    PERFORM check_tender(usr, tender_id);
    IF nm IS NOT NULL THEN
        UPDATE tender SET name=nm WHERE id=tender_id;
    END IF;
    IF des IS NOT NULL THEN
        UPDATE tender SET description=des WHERE id=tender_id;
    END IF;
    IF tp IS NOT NULL THEN
        UPDATE tender SET type=tp WHERE id=tender_id;
    END IF;
    IF st IS NOT NULL THEN
        UPDATE tender SET status=st WHERE id=tender_id;
    END IF;
    SELECT * INTO t FROM tender WHERE id=tender_id;
    UPDATE tender SET version=version+1 WHERE id=tender_id;
    RETURN t;
END
$$ LANGUAGE plpgsql;

create function create_tender(usr character varying, nm character varying, des text, tp service_type, org_id uuid) returns tender
    language plpgsql
as
$$
DECLARE
    uid UUID = get_user(usr);
    t tender;
BEGIN
    IF NOT EXISTS(SELECT 1 FROM employee e JOIN organization_responsible o ON e.id=o.user_id WHERE e.id = uid AND o.organization_id=org_id) THEN
        RAISE 'no rights' USING ERRCODE='ACCESS_DENIED';
END IF;
INSERT INTO tender (name, description, type, organization_id, created_by) VALUES (nm, des, tp, org_id, usr) RETURNING * INTO t;
RETURN t;
END;
$$;