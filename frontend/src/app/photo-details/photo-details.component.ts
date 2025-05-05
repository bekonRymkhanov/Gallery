import { NgIf, NgForOf, NgClass, CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { GalleryService } from '../gallery.service';
import { Photo, Rating,Comment } from '../models';

@Component({
  selector: 'app-photo-details',
  standalone: true,
  imports: [NgIf,
            NgForOf,
            NgClass,
            CommonModule,
            RouterLink,
            FormsModule,
            ReactiveFormsModule
  ],
  templateUrl: './photo-details.component.html',
  styleUrl: './photo-details.component.css'
})
export class PhotoDetailsComponent {
  photoId: number = 0;
  photo: Photo | undefined;
  comments: Comment[] = [];
  averageRating: number = 5;
  totalRatings: number = 0;
  isLoading: boolean = true;
  commentForm: FormGroup;
  currentUserRating: number = 0;
  isLiked: boolean = false;
  errorMessage: string = '';
  successMessage: string = '';
  isAuthenticated: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private galleryService: GalleryService,
    private fb: FormBuilder
  ) {
    this.commentForm = this.fb.group({
      content: ['', [Validators.required, Validators.minLength(2)]]
    });
  }

  ngOnInit(): void {
    this.isAuthenticated = sessionStorage.getItem('currentUser') !== null;
    
    this.route.params.subscribe(params => {
      this.photoId = params['photoid'];
      this.galleryService.getPhoto(this.photoId).subscribe(request => {
        this.photo = request.photo;
        this.loadComments();
        this.loadRatings();
        if (this.isAuthenticated) {
          this.checkUserLike();
          this.getUserRating();
        }
      });
    });
  }

  
  loadComments(): void {
    this.galleryService.getPhotoComments(this.photoId).subscribe({
      next: (response) => {
        this.comments = response.comments;
      },
      error: (error) => {
        this.errorMessage = 'Failed to load comments';
      }
    });
  }

  loadRatings(): void {

    this.galleryService.getPhotoRating(this.photoId).subscribe(request => {
        this.averageRating = request.average_score;
        this.totalRatings = request.count;
        this.isLoading = false;

      },
    );
  }

  getUserRating(): void {
    this.galleryService.getPhotoRatingByUser(this.photoId).subscribe(request => {
        this.currentUserRating = request.rating.score;
      },

    );
  }

  checkUserLike(): void {
    this.galleryService.checkLike(this.photoId).subscribe({
      next: (response) => {
        this.isLiked = response.liked;
      }
    });
  }

  submitComment(): void {
    if (this.commentForm.invalid) {
      return;
    }

    if (!this.isAuthenticated) {
      this.errorMessage = 'You must be logged in to comment';
      return;
    }

    const newComment = {
      photo_id: Number(this.photoId),
      content: this.commentForm.value.content
    };

    this.galleryService.postComment(newComment).subscribe({
      next: (comment) => {
        this.commentForm.reset();
        this.comments.unshift(comment);
        this.successMessage = 'Comment added successfully';
        setTimeout(() => this.successMessage = '', 3000);
      },
      error: (error) => {
        this.errorMessage = 'Failed to add comment';
      }
    });
  }

  ratePhoto(score: number): void {
    if (!this.isAuthenticated) {
      this.errorMessage = 'You must be logged in to rate';
      return;
    }

    const newRating = {
      photo_id: Number(this.photoId),
      score: score,
    };

    this.galleryService.postRating(newRating).subscribe({
      next: (rating) => {
        this.currentUserRating = score;
      
        this.successMessage = 'Rating submitted successfully';
        setTimeout(() => this.successMessage = '', 3000);
      },
      error: (error) => {
        this.errorMessage = 'Failed to submit rating';
      }
    });
  }

  toggleLike(): void {
    if (!this.isAuthenticated) {
      this.errorMessage = 'You must be logged in to like photos';
      return;
    }

    if (this.isLiked) {
      this.galleryService.unlikePhoto(this.photoId).subscribe({
        next: () => {
          this.isLiked = false;
          if (this.photo) {
            this.photo.likes--;
          }
        },
        error: (error) => {
          this.errorMessage = 'Failed to unlike photo';
        }
      });
    } else {
      this.galleryService.likePhoto(this.photoId).subscribe({
        next: () => {
          this.isLiked = true;
          if (this.photo) {
            this.photo.likes++;
          }
        },
        error: (error) => {
          this.errorMessage = 'Failed to like photo';
        }
      });
    }
  }

  deleteComment(commentId: number): void {
    if (confirm('Are you sure you want to delete this comment?')) {
      this.galleryService.deleteComment(commentId).subscribe({
        next: () => {
          this.comments = this.comments.filter(c => c.id !== commentId);
          this.successMessage = 'Comment deleted successfully';
          setTimeout(() => this.successMessage = '', 3000);
        },
        error: (error) => {
          this.errorMessage = 'Failed to delete comment';
        }
      });
    }
  }

  formatDate(date: Date): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}
