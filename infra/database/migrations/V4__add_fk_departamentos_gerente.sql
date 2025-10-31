ALTER TABLE departamentos 
ADD CONSTRAINT fk_departamentos_gerente 
FOREIGN KEY (gerente_id) REFERENCES colaboradores(id) DEFERRABLE INITIALLY DEFERRED;

