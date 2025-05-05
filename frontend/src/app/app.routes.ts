import { Router, Routes } from '@angular/router';
import {HomeComponent} from "./home/home.component";
import {AboutComponent} from "./about/about.component";
import {NotFoundComponent} from "./not-found/not-found.component";
import { inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ProfileComponent } from './profile/profile.component';
import { PhotosComponent } from './photos/photos.component';
import { PhotoDetailsComponent } from './photo-details/photo-details.component';
import { LikedPhotosComponent } from './liked-photos/liked-photos.component';


export const routes: Routes = [
    { path:"",redirectTo:"home",pathMatch:"full" },
    { path:"home",component:HomeComponent,title:"Home page"},
    { path:"about",component:AboutComponent,title:"About page" },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'profile', component: ProfileComponent },
    { path: 'photos/:photoid', component: PhotoDetailsComponent ,title:"Photos details"},
    { path: 'photos', component: PhotosComponent ,title:"Photos page"},
    { path: 'likedPhotos', component: LikedPhotosComponent ,title:"Photos page"},


    { path:"**",component:NotFoundComponent,title:"404 - not found" }
];
