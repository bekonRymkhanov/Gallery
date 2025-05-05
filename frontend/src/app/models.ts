


export interface Comment {
  id: number;
  photo_id: number;
  user_id: number;
  content: string;
  created_at: Date;
  version: number;
}

export interface Rating {
  id: number;
  photo_id: number;
  user_id:number;
  score: number;
  created_at: Date;
  version: number;
}
export interface RatingResponse {
  rating: Rating;
}

export interface PhotoRatingsResponse {
  average_score: number;
  count: number;
  ratings: Rating[];
}

export interface Metadata {
  CurrentPage: number;
  PageSize: number;
  FirstPage: number;
  LastPage: number;
  TotalRecords: number;
}



export interface User {
  id?: number;
  name: string;
  email: string;
  password?: string;
  is_admin: boolean;
  activated?: boolean;
  created_at?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ActivationResponse {
  user: User;
}

export interface AuthenticationResponse {
  authentication_token: {
    token: string;
    expiry: string;
  }
  user:User
}

export interface Photo {
  id: number;
  title: string;
  description: string;
  author: string;
  category: string;
  tags: string;
  width: number;
  height: number;
  url: string;
  thumbnail_url: string;
  source: string;
  download_count: number;
  likes: number;
  version: number;
}

export interface PhotoFilters {
  title?: string;
  author?: string;
  category?: string;
  tags?: string;
  page?: number;
  page_size?: number;
  sort?: string;
}


export interface Like {
  id: number;
  user_id: number;
  photo_id: number;
  created_at: Date;
  version: number;
}
export interface PhotoResponse {
  photos: Photo[];
}


export interface requestPhotoDetail {
  photo: Photo
}