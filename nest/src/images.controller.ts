import {
  Controller,
  Post,
  Get,
  Put,
  Delete,
  Body,
  Param,
  Res,
} from '@nestjs/common';
import { ImagesService } from './images.service';
import { Response } from 'express';

@Controller('')
export class ImagesController {
  constructor(private readonly imagesService: ImagesService) {}

  @Post('/create')
  async createImages(@Res() res: Response) {
    try {
      const result = await this.imagesService.createImages();
      return res.json(result);
    } catch (error) {
      return res.status(500).json({ error: error.message });
    }
  }

  @Get('/read')
  async getAllImages(@Res() res: Response) {
    try {
      const result = await this.imagesService.getAllImages();
      return res.json(result);
    } catch (error) {
      return res.status(500).json({ error: error.message });
    }
  }

  @Put('/update/:id')
  async updateImage(
    @Param('id') id: string,
    @Body() updateData: any,
    @Res() res: Response,
  ) {
    try {
      const result = await this.imagesService.updateImage(
        Number(id),
        updateData,
      );
      return res.json(result);
    } catch (error) {
      console.log(error);

      return res.status(500).json({ error: error.message });
    }
  }

  @Delete('/delete/:id')
  async deleteImage(@Param('id') id: string, @Res() res: Response) {
    try {
      const result = await this.imagesService.deleteImage(Number(id));
      if (result) {
        return res.json(result);
      } else {
        return res
          .status(404)
          .json({ message: 'No se encontró la imagen a eliminar' });
      }
    } catch (error) {
      return res.status(500).json({ error: error.message });
    }
  }
}
