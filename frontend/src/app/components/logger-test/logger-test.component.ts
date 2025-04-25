import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-logger-test',
  template: `
    <div class="container">
      <h2>Kafka Logger Test</h2>
      <button (click)="testLogging()" class="btn btn-primary">Test Logging</button>
      <div *ngIf="response" class="mt-3">
        <p>Response: {{ response | json }}</p>
      </div>
    </div>
  `,
  styles: [`
    .container {
      padding: 20px;
      max-width: 600px;
      margin: 0 auto;
    }
    .btn {
      padding: 10px 20px;
      border-radius: 5px;
      cursor: pointer;
    }
    .btn-primary {
      background-color: #007bff;
      color: white;
      border: none;
    }
  `]
})
export class LoggerTestComponent {
  response: any;

  constructor(private http: HttpClient) {}

  testLogging() {
    this.http.get('http://localhost:3001/logger/test').subscribe(
      (response) => {
        this.response = response;
        console.log('Logging test successful:', response);
      },
      (error) => {
        console.error('Error testing logging:', error);
        this.response = { error: error.message };
      }
    );
  }
} 