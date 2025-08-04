import { Component, Input } from '@angular/core';
import { LabService } from '../../services/lab.service';
import { RatingModule } from 'primeng/rating';
import { FloatLabelModule } from 'primeng/floatlabel';
import { ButtonModule } from 'primeng/button';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-session-feedback',
  imports: [RatingModule, FloatLabelModule, ButtonModule, FormsModule],
  templateUrl: './session-feedback.component.html',
  styleUrl: './session-feedback.component.css'
})
export class SessionFeedbackComponent {
  @Input() sessionId: string;
  @Input() rating: number = 0;
  feedback: string = '';

  constructor(private labService: LabService) {
    this.sessionId = '';
  }

  sendFeedback() {
    this.labService.sendFeedback(this.sessionId, this.rating, this.feedback).subscribe((response) => {
      console.log(response);
    });
  }
}
