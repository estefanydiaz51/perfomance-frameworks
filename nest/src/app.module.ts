import { Module } from '@nestjs/common';
import { ImagesModule } from './images.module';

@Module({
  imports: [ImagesModule],
  controllers: [],
  providers: [],
})
export class AppModule {}
