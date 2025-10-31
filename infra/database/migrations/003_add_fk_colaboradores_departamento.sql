ALTER TABLE colaboradores 
ADD CONSTRAINT fk_colaboradores_departamento 
FOREIGN KEY (departamento_id) REFERENCES departamentos(id) DEFERRABLE INITIALLY DEFERRED


