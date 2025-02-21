import {
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  ElementRef,
  ViewChild,
} from '@angular/core';
import { defineElement } from '@lordicon/element';
import lottie from 'lottie-web';

@Component({
  selector: 'app-home',
  imports: [],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class HomeComponent {
  @ViewChild('musicIcon') musicIconElement!: ElementRef;
  @ViewChild('filmIcon') filmIconElement!: ElementRef;

  constructor() {
    defineElement(lottie.loadAnimation);
  }

  playAnimation(el: ElementRef) {
    const element = el.nativeElement;
    element.playerInstance.playFromBeginning();
  }
}
