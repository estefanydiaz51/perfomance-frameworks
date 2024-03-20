const { Pool } = require('pg');

const pool = new Pool({
  user: 'usuario',
  host: 'localhost',
  database: 'crud_db',
  password: '5131',
  port: 5433,
});

module.exports = pool;
