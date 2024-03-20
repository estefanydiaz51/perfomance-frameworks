import { Module } from '@nestjs/common';
import { Pool } from 'pg';

const pool: Pool = new Pool({
  user: 'usuario',
  host: 'localhost',
  database: 'crud_db',
  password: '5131',
  port: 5433,
});

@Module({
  providers: [
    {
      provide: 'DATABASE_POOL',
      useValue: pool,
    },
  ],
  exports: ['DATABASE_POOL'],
})
export class DatabaseModule {}
