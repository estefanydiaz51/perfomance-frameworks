import json
import time
import psycopg2
from psycopg2.extras import RealDictCursor
from aiofile import async_open
from flask import Flask, request, jsonify

app = Flask(__name__)

# Configuración de la conexión a la base de datos.
conn = psycopg2.connect(host='localhost', database='crud_db',
                        user='usuario', password='5131', port='5433')
conn.autocommit = True

# CREATE


@app.route('/create', methods=['POST'])
async def create():
    try:
        async with async_open('dummy.json', 'r') as afp:
            json_data = await afp.read()
        images = json.loads(json_data)

        with conn.cursor() as cur:
            cur.execute('SELECT 1 FROM images LIMIT 1')
            exist_images = cur.fetchall()
            if exist_images:
                cur.execute('DELETE FROM images')
            start_time = time.perf_counter()
            for image in images:
                cur.execute(
                    'INSERT INTO images (albumId, id, title, url, thumbnailUrl) VALUES (%s, %s, %s, %s, %s) ON CONFLICT (id) DO NOTHING RETURNING *',
                    (image['albumId'], image['id'], image['title'],
                     image['url'], image['thumbnailUrl'])
                )
            end_time = time.perf_counter()
            time_taken = end_time - start_time
            return jsonify({
                'message': 'Datos cargados correctamente',
                'time': f'{time_taken} ms',
                'dataQuantity': len(images)
            })
    except Exception as error:
        print(f'Error al procesar la solicitud: {error}')
        return {'error': 'Error al cargar los datos'}, 500

# READ


@app.route('/read', methods=['GET'])
def read():
    try:
        start_time = time.perf_counter()
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute('SELECT * FROM images')
            all_images = cur.fetchall()
        end_time = time.perf_counter()
        time_taken = end_time - start_time
        return jsonify({
            'data': sorted(all_images, key=lambda x: x['id']),
            'dataQuantity': len(all_images),
            'time': f'{time_taken} ms',
        })
    except Exception as error:
        print(f'Error al procesar la solicitud: {error}')
        return {'error': 'Error al leer los datos'}, 500

# UPDATE


@app.route('/update/<int:id>', methods=['PUT'])
def update(id):
    try:
        data = request.json
        start_time = time.perf_counter()
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute('UPDATE images SET albumId = %s, title = %s, url = %s, thumbnailUrl = %s WHERE id = %s RETURNING *',
                        (data['albumId'], data['title'], data['url'], data['thumbnailUrl'], id))
            updated_image = cur.fetchone()
        end_time = time.perf_counter()
        time_taken = end_time - start_time
        return jsonify({'data': updated_image, 'time': f'{time_taken} ms'})
    except Exception as error:
        print(f'Error al procesar la solicitud: {error}')
        return {'error': 'Error al actualizar los datos'}, 500

# DELETE


@app.route('/delete/<int:id>', methods=['DELETE'])
def delete(id):
    try:
        start_time = time.perf_counter()
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute('DELETE FROM images WHERE id = %s RETURNING *', (id,))
            deleted_image = cur.fetchone()
        end_time = time.perf_counter()
        time_taken = end_time - start_time
        return jsonify({'data': deleted_image, 'time': f'{time_taken} ms'})
    except Exception as error:
        print(f'Error al procesar la solicitud: {error}')
        return {'error': 'Error al borrar los datos'}, 500


if __name__ == '__main__':
    app.run(debug=True, port=3000)
