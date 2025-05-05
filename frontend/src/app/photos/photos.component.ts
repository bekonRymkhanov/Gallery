import { Component, OnInit } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, FormsModule } from '@angular/forms';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { GalleryService } from '../gallery.service';
import { Photo, PhotoFilters, Metadata } from '../models';
import { NgIf, NgForOf } from '@angular/common';

@Component({
  selector: 'app-photos',
  standalone: true,
  imports: [
    NgIf,
    FormsModule,
    RouterLink,
    NgForOf
  ],
  templateUrl: './photos.component.html',
  styleUrl: './photos.component.css'
})
export class PhotosComponent implements OnInit {
  photos: Photo[] = [];
  metadata!: Metadata;
  loaded = false;

  filters: PhotoFilters = {
    title: '',
    author: '',
    category: '',
    tags: '',
  };

  constructor(private httpService:GalleryService) {}

  ngOnInit(): void {
    this.loadPhotos();
  }

  loadPhotos(): void {
    this.loaded = false;
    this.httpService.getPhotos(this.filters).subscribe(response => {
      this.photos= response.photos;
      this.metadata = response.metadata;
      this.loaded = true;
    });

  }

  onSearch(): void {
    this.filters = { ...this.filters, page: 1 };
    
    this.loadPhotos();
  }

  onClear(): void {
    this.filters = { title: '', author: '', category: '', tags: '', page: 1 };
    this.loadPhotos();
  }
  NextPage(): void {
    if (this.metadata.CurrentPage >= this.metadata.LastPage) {
      return;
    }
    this.filters = { ...this.filters, page: this.metadata.CurrentPage + 1 };
    this.loadPhotos();
  }
  PreviousPage(): void {
    if (this.metadata.CurrentPage <= 1) {
      return;
    }
    this.filters = { ...this.filters, page: this.metadata.CurrentPage - 1 };
    this.loadPhotos();
  }

}