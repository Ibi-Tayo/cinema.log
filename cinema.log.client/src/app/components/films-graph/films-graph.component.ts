import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  ViewChild,
} from '@angular/core';
import { DataSet } from 'vis-data';
import { Data, Edge, Network, Options } from 'vis-network';

@Component({
  selector: 'app-films-graph',
  imports: [],
  templateUrl: './films-graph.component.html',
  styleUrl: './films-graph.component.scss',
})
export class FilmsGraphComponent implements AfterViewInit, OnDestroy {
  @ViewChild('networkContainer') networkContainer!: ElementRef;

  private network: Network | undefined;

  ngAfterViewInit(): void {
    // 1. Define nodes
    const nodes = new DataSet([
      { id: 1, label: 'Film A (Seed Film)' },
      { id: 2, label: 'Film B' },
      { id: 3, label: 'Film C' },
      { id: 4, label: 'Film D' },
      { id: 5, label: 'Film E' },
    ]);

    // 2. Define edges
    const edges: DataSet<Edge> = new DataSet([
      { id: 1, from: 1, to: 2 },
      { id: 2, from: 1, to: 3 },
      { id: 3, from: 1, to: 4 },
      { id: 4, from: 1, to: 5 },
    ]);

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
