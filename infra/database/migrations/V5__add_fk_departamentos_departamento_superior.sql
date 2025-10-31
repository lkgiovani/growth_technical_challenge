ALTER TABLE departamentos 
ADD CONSTRAINT fk_departamentos_departamento_superior 
FOREIGN KEY (departamento_superior_id) REFERENCES departamentos(id);

