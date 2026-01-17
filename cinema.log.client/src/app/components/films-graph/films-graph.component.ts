import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  ViewChild,
} from '@angular/core';
import { DataSet } from 'vis-data';
import { Data, Edge, Network, Options } from 'vis-network';
import { GraphService } from '../../services/graph.service';

@Component({
  selector: 'app-films-graph',
  imports: [],
  templateUrl: './films-graph.component.html',
  styleUrl: './films-graph.component.scss',
})
export class FilmsGraphComponent implements AfterViewInit, OnDestroy {
  @ViewChild('networkContainer') networkContainer!: ElementRef;

  private network: Network | undefined;

  constructor(private graphService: GraphService) {}

  ngAfterViewInit(): void {
    this.loadGraphData();
  }

  private loadGraphData(): void {
    this.graphService.getUserGraph().subscribe({
      next: (graphData) => {
        this.renderGraph(graphData);
      },
      error: (error) => {
        console.error('Error loading graph data:', error);
        // Optionally show empty graph or error message
        this.renderGraph({ nodes: [], edges: [] });
      },
    });
  }

  private renderGraph(graphData: {
    nodes: Array<{ externalFilmId: number; title: string }>;
    edges: Array<{ fromFilmId: number; toFilmId: number }>;
  }): void {
    // 1. Transform nodes to vis-network format
    const nodes = new DataSet(
      graphData.nodes.map((node) => ({
        id: node.externalFilmId,
        label: node.title,
      }))
    );

    // 2. Transform edges to vis-network format
    const edges: DataSet<Edge> = new DataSet(
      graphData.edges.map((edge, index) => ({
        id: index + 1,
        from: edge.fromFilmId,
        to: edge.toFilmId,
      }))
    );

    // 3. Config
    const options: Options = {
      physics: {
        enabled: true,
        barnesHut: {
          gravitationalConstant: -2000,
          centralGravity: 0.3,
          springLength: 95,
        },
      },
      nodes: {
        shape: 'dot',
        size: 16,
        font: { size: 14, color: '#fbfbfbff' },
        borderWidth: 2,
      },
      edges: {
        width: 2,
        color: { color: '#848484' },
      },
    };
    // 4. Create network
    const data: Data = { nodes, edges };
    this.network = new Network(
      this.networkContainer.nativeElement,
      data,
      options
    );
  }

  ngOnDestroy(): void {
    if (this.network) {
      this.network?.destroy();
    }
  }
}
