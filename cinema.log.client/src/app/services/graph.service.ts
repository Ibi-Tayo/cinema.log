import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

export interface FilmGraphNode {
  userId: string;
  externalFilmId: number;
  title: string;
}

export interface FilmGraphEdge {
  userId: string;
  edgeId: string;
  fromFilmId: number;
  toFilmId: number;
}

export interface UserGraph {
  nodes: FilmGraphNode[];
  edges: FilmGraphEdge[];
}

@Injectable({
  providedIn: 'root',
})
export class GraphService {
  private http = inject(HttpClient);


  getUserGraph(): Observable<UserGraph> {
    return this.http.get<UserGraph>(`${import.meta.env.NG_APP_API_URL}/graph`, {
      withCredentials: true,
    });
  }
}
