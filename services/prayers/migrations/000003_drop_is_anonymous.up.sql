-- is_anonymous decía lo mismo que author_name IS NULL. Dos columnas para un
-- hecho es una que puede contradecir a la otra: la fila podía afirmar que la
-- petición era anónima y llevar el nombre escrito al lado.
ALTER TABLE prayers DROP COLUMN is_anonymous;
