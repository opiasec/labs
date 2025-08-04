import { Component, ElementRef, AfterViewInit, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import * as Diff2Html from 'diff2html';


@Component({
  selector: 'app-files-diff-viewer',
  templateUrl: './files-diff-viewer.component.html',
  styleUrls: ['./files-diff-viewer.component.css'],
  standalone: true,
  imports: [CommonModule]
})
export class FilesDiffViewerComponent implements AfterViewInit {
  @Input() diff: string = '';

  constructor(private el: ElementRef) {}

  ngAfterViewInit() {
    const target = this.el.nativeElement.querySelector('#diffContainer');
    if (target && this.diff) {
      const html = Diff2Html.html(this.diff, {
        drawFileList: true,
        matching: 'lines',
        outputFormat: 'side-by-side',
      });
      const darkHtml = html.replace(/d2h-light-color-scheme/g, 'd2h-dark-color-scheme');

      target.innerHTML = darkHtml;
    }
  }
}
