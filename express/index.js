const express = require('express');
const fs = require('fs/promises');
const pool = require('./db');

const app = express();
const PORT = 3000;

app.use(express.json());

// CREATE
app.post('/create', async (req, res) => {
  try {
    const jsonData = await fs.readFile('dummy.json', 'utf8');
    const images = JSON.parse(jsonData);
    const existImages = await pool.query('SELECT 1 FROM images LIMIT 1');
    if (existImages.rowCount > 0) {
      await pool.query('DELETE FROM images');
    }
    const startTime = Date.now();
    for await (const image of images) {
      await pool.query(
        'INSERT INTO images (albumId, id, title, url, thumbnailUrl) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING RETURNING *',
        [image.albumId, image.id, image.title, image.url, image.thumbnailUrl]
      );
    }

    const endTime = Date.now();
    const timeTaken = endTime - startTime;

    res.json({
      message: 'Datos cargados correctamente',
      time: `${timeTaken} ms`,
      dataQuantity: images.length,
    });
  } catch (error) {
    console.error('Error al procesar la solicitud:', error);
    res.status(500).send({ error: 'Error al cargar los datos' });
  }
});

// READ
app.get('/read', async (req, res) => {
  try {
    const startTime = Date.now();
    const allImages = await pool.query('SELECT * FROM images');
    const endTime = Date.now();
    const timeTaken = endTime - startTime;
    res.json({
      data: allImages.rows.sort((a, b) => a.id - b.id),
      dataQuantity: allImages.rows.length,
      time: `${timeTaken} ms`,
    });
  } catch (error) {
    console.error('Error al procesar la solicitud:', error);
    res.status(500).send({ error: 'Error al leer los datos' });
  }
});

// UPDATE
app.put('/update/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const { albumId, title, url, thumbnailUrl } = req.body;
    const startTime = Date.now();
    const uploadImage = await pool.query(
      'UPDATE images SET albumId = $1, title = $2, url = $3, thumbnailUrl = $4 WHERE id = $5 RETURNING *',
      [albumId, title, url, thumbnailUrl, id]
    );
    const endTime = Date.now();
    const timeTaken = endTime - startTime;
    res.json({ data: uploadImage.rows[0], time: `${timeTaken} ms` });
  } catch (error) {
    console.error('Error al procesar la solicitud:', error);
    res.status(500).send({ error: 'Error al actualizar los datos' });
  }
});

// DELETE
app.delete('/delete/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const startTime = Date.now();
    const deleteImage = await pool.query(
      'DELETE FROM images WHERE id = $1 RETURNING *',
      [id]
    );
    const endTime = Date.now();
    const timeTaken = endTime - startTime;
    res.json({ data: deleteImage.rows[0], time: `${timeTaken} ms` });
  } catch (error) {
    console.error('Error al procesar la solicitud:', error);
    res.status(500).send({ error: 'Error al borrar los datos' });
  }
});

app.listen(PORT, () => {
  console.log(`Servidor corriendo en http://52.23.103.132:${PORT}`);
});
