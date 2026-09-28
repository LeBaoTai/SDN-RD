CREATE TABLE public.devices (
    id uuid NOT NULL PRIMARY KEY,
    name varchar(50) NOT NULL,
    address inet NOT NULL, -- Định nghĩa kiểu IP của Postgres
    port integer NOT NULL,

    CONSTRAINT check_valid_port CHECK (port BETWEEN 0 AND 65535),

    CONSTRAINT unique_device_endpoint UNIQUE (address, port)
);

CREATE INDEX idx_devices_address ON public.devices (address);

ALTER TABLE public.devices 
ALTER COLUMN id SET DEFAULT gen_random_uuid();
