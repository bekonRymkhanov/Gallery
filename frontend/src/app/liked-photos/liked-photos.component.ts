import { Component, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { CommonModule, NgIf, NgForOf } from '@angular/common';
import { GalleryService } from '../gallery.service';
import { Photo, Like, PhotoResponse } from '../models';

@Component({
  selector: 'app-liked-photos',
  standalone: true,
  imports: [
    CommonModule,
    NgIf,
    NgForOf,
    RouterLink
  ],
  templateUrl: './liked-photos.component.html',
  styleUrl: './liked-photos.component.css'
})
export class LikedPhotosComponent implements OnInit {
  likedPhotos: Photo[] = [];
  loaded = false;
  error: string | null = null;

  constructor(private galleryService: GalleryService) {}

  ngOnInit(): void {
    this.loadLikedPhotos();
  }

  loadLikedPhotos(): void {
    this.loaded = false;
    this.error = null;
    
    this.galleryService.getLikesOfUser().subscribe(response => {
      console.log(response)

      this.likedPhotos = response.photos;
      this.loaded = true;
      },

    );
  }

  unlikePhoto(photoId: number): void {
    this.galleryService.unlikePhoto(photoId).subscribe({
      });
    this.likedPhotos = this.likedPhotos.filter(photo => photo.id !== photoId);

  }
}