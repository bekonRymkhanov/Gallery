import {Component, OnInit} from '@angular/core';
import {RouterLink, RouterOutlet} from '@angular/router';
import { CommonModule } from '@angular/common';
import {FormsModule} from "@angular/forms";
import {GalleryService} from "./gallery.service";
import {HomeComponent} from "./home/home.component";
import {AboutComponent} from "./about/about.component";
import { NotFoundComponent } from './not-found/not-found.component';
import { Router } from '@angular/router';
import { AuthService } from './auth.service';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ProfileComponent } from './profile/profile.component';
import { PhotosComponent } from './photos/photos.component';
import { PhotoDetailsComponent } from './photo-details/photo-details.component';
import { LikedPhotosComponent } from './liked-photos/liked-photos.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [PhotosComponent,PhotoDetailsComponent,LikedPhotosComponent,RouterOutlet,ProfileComponent,LoginComponent,NotFoundComponent,HomeComponent,RegisterComponent, CommonModule,  RouterLink, FormsModule,HomeComponent,AboutComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit{
  title= 'Gallery';
  isLoggedIn:boolean=false;
  constructor(private httpService:GalleryService,private authService:AuthService) {
    this.authService=authService;
  }

  ngOnInit(): void {
    this.isLoggedIn=this.authService.isLoggedIn;
  }

  protected readonly localStorage = localStorage;
  logout() {
    this.authService.logout();
    this.isLoggedIn = false;
  }
}
