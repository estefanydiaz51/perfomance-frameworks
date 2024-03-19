import { Injectable, Inject } from '@nestjs/common';
import { Pool } from 'pg';
import * as fs from 'fs';

@Injectable()
export class ImagesService {
  constructor(@Inject('DATABASE_POOL') private pool: Pool) {}

  async createImages(): Promise<any> {
    try {
      const jsonData = await fs.readFileSync('dummy.json', 'utf8');
      const images = JSON.parse(jsonData);
      const existImages = await this.pool.query('SELECT 1 FROM images LIMIT 1');
      if (existImages.rowCount > 0) {
        await this.pool.query('DELETE FROM images');
      }
      const startTime = Date.now();
      for (const image of images) {
        await this.pool.query(
          'INSERT INTO images (albumId, id, title, url, thumbnailUrl) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING RETURNING *',
          [image.albumId, image.id, image.title, image.url, image.thumbnailUrl],
        );
      }
      const endTime = Date.now();
      const timeTaken = endTime - startTime;
      return {
        message: 'Datos cargados correctamente',
        time: `${timeTaken} ms`,
        dataQuantity: images.length,
      };
    } catch (error) {
      console.log(error);

      throw new Error('Error al cargar los datos');
    }
  }

  async getAllImages(): Promise<any> {
    try {
      const startTime = Date.now();

      const allImages = await this.pool.query('SELECT * FROM images');
      const endTime = Date.now();
      const timeTaken = endTime - startTime;
      return {
        data: allImages.rows.sort((a, b) => a.id - b.id),
        dataQuantity: allImages.rows.length,
        time: `${timeTaken} ms`,
      };
    } catch (error) {
      throw new Error('Error al leer los datos');
    }
  }

  async updateImage(id: number, updateData: any): Promise<any> {
    try {
      const { albumId, title, url, thumbnailUrl } = updateData;
      const startTime = Date.now();

      const updatedImage = await this.pool.query(
        'UPDATE images SET albumId = $1, title = $2, url = $3, thumbnailUrl = $4 WHERE id = $5 RETURNING *',
        [albumId, title, url, thumbnailUrl, id],
      );
      const endTime = Date.now();
      const timeTaken = endTime - startTime;
      return {
        data: updatedImage.rows[0],
        time: `${timeTaken} ms`,
      };
    } catch (error) {
      throw new Error('Error al actualizar los datos');
    }
  }

  async deleteImage(id: number): Promise<any> {
    try {
      const startTime = Date.now();

      const deleteImage = await this.pool.query(
        'DELETE FROM images WHERE id = $1 RETURNING *',
        [id],
      );
      const endTime = Date.now();
      const timeTaken = endTime - startTime;

      return {
        data: deleteImage.rows[0],
        time: `${timeTaken} ms`,
      };
    } catch (error) {
      throw new Error('Error al borrar los datos');
    }
  }
}
