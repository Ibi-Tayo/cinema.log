import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

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
  constructor(private http: HttpClient) {}

  getUserGraph(): Observable<UserGraph> {
    return this.http.get<UserGraph>(`${environment.apiUrl}/graph`, {
      withCredentials: true,
    });
  }
}
