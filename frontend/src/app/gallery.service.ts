import { Injectable } from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {User,Comment,Rating,Metadata, AuthenticationResponse,Photo,Like,PhotoFilters, requestPhotoDetail, PhotoRatingsResponse, RatingResponse, PhotoResponse} from "./models";
import { Observable } from 'rxjs/internal/Observable';

@Injectable({
  providedIn: 'root'
})
export class GalleryService {
  BACKEND_URL='http://4.247.165.79:8080';

  constructor(private client:HttpClient) { }

  private getAuthHeaders(): HttpHeaders {
    const userData = sessionStorage.getItem('currentUser');
    if (userData) {
      const user: { 
        token: string;
        expiry: string;
        user:User}= JSON.parse(userData);
      return new HttpHeaders({
        'Authorization': `Bearer ${user.token}`
      });
    }
    return new HttpHeaders();
  }
  getPhotos(filters:PhotoFilters){
    return this.client.get<{ photos: Photo[], metadata: Metadata }>(`${this.BACKEND_URL}/photos/`, { params: { ...filters } })
  }
  getPhoto(id:number){
    const headers = this.getAuthHeaders();
    return this.client.get<requestPhotoDetail>(`${this.BACKEND_URL}/photos/${id}`, { headers })
  }
  deletePhoto(id:number){
    const headers = this.getAuthHeaders();

    return this.client.delete(`${this.BACKEND_URL}/photos/${id}/`, { headers })
  }
  postPhoto(newPhoto:Photo){
    const headers = this.getAuthHeaders();
    return this.client.post<Photo>(`${this.BACKEND_URL}/photo/`,newPhoto, { headers })
  }
  putPhoto(newphoto:Photo){
    const headers = this.getAuthHeaders();
    return this.client.patch<Photo>(`${this.BACKEND_URL}/photos/${newphoto.id}/`,newphoto, { headers })
  }
  getPhotoComments(id:number){
    const headers = this.getAuthHeaders();

    return this.client.get<{comments: Comment[], metadata: Metadata}>(`${this.BACKEND_URL}/photos/${id}/comments/`, { headers })
  }
  getPhotoRating(id:number){
    const headers = this.getAuthHeaders();

    return this.client.get<PhotoRatingsResponse>(`${this.BACKEND_URL}/photos/${id}/ratings/`, { headers })
  }


  getComment(id:number){
    const headers = this.getAuthHeaders();

    return this.client.get<Comment>(`${this.BACKEND_URL}/comments/${id}/`, { headers })
  }
  deleteComment(id:number){
    const headers = this.getAuthHeaders();

    return this.client.delete(`${this.BACKEND_URL}/comments/${id}/`, { headers })
  }
  postComment(newComment:{
    photo_id: number;
    content: string;
  }){
    const headers = this.getAuthHeaders();

    return this.client.post<Comment>(`${this.BACKEND_URL}/comments/`,newComment, { headers })
  }
  putComment(newComment:Comment){
    const headers = this.getAuthHeaders();

    return this.client.patch<Comment>(`${this.BACKEND_URL}/comments/${newComment.id}/`,newComment, { headers })
  }



  getRating(id:number){
    const headers = this.getAuthHeaders();

    return this.client.get<Rating>(`${this.BACKEND_URL}/ratings/${id}/`, { headers })
  }
  deleteRating(id:number){
    const headers = this.getAuthHeaders();

    return this.client.delete(`${this.BACKEND_URL}/ratings/${id}/`, { headers })
  }
  postRating(newRating:{
    photo_id:number,
    score: number
  }){
    const headers = this.getAuthHeaders();

    return this.client.post<Rating>(`${this.BACKEND_URL}/ratings/`,newRating, { headers })
  }
  putRating(newRating:Rating){
    const headers = this.getAuthHeaders();

    return this.client.patch<Rating>(`${this.BACKEND_URL}/ratings/${newRating.id}/`,newRating , { headers });
  }
  getPhotoRatingByUser(photoID:number){
    const headers = this.getAuthHeaders();
    return this.client.get<RatingResponse>(`${this.BACKEND_URL}/photo/${photoID}/rating/`, { headers });
  }



  getLikesOfUser() {
    const headers = this.getAuthHeaders();
    return this.client.get<PhotoResponse>(`${this.BACKEND_URL}/users/likes/`, { headers });
  }
  checkLike(photoID: number) {
    const headers = this.getAuthHeaders();
    return this.client.get<{liked: boolean;}>(`${this.BACKEND_URL}/photos/${photoID}/likes/check/`, { headers });
  }
  likePhoto(photoID: number) {
    const headers = this.getAuthHeaders();
    return this.client.post<Like>(`${this.BACKEND_URL}/photos/${photoID}/likes/`, {}, { headers });
  }
  unlikePhoto(photoID: number) {
    const headers = this.getAuthHeaders();
    return this.client.delete(`${this.BACKEND_URL}/photos/${photoID}/likes/`, { headers });
  }
}