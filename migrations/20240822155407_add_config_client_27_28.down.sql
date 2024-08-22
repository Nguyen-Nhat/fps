DELETE FROM fps_client WHERE client_id in (27, 28);
DELETE FROM config_mapping WHERE client_id in (27, 28);
DELETE FROM config_task WHERE config_mapping_id in (27, 28);
