import { APP_INITIALIZER, ApplicationConfig, inject, InjectionToken, provideAppInitializer } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { provideClientHydration } from '@angular/platform-browser';
import { providePrimeNG } from 'primeng/config';
import lara from '@primeng/themes/lara';
import Aura from '@primeuix/themes/aura';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor } from './interceptor/auth.interceptor';
import { MessageService } from 'primeng/api';
import { CLIPBOARD_OPTIONS, ClipboardButtonComponent, MARKED_OPTIONS, provideMarkdown } from 'ngx-markdown';
import 'prismjs';
import 'prismjs/components/prism-typescript.min.js';
import 'prismjs/components/prism-javascript.min.js';
import 'prismjs/components/prism-css.min.js';
import 'prismjs/components/prism-markup.min.js';
import 'prismjs/components/prism-bash.min.js';
import 'prismjs/components/prism-json.min.js';
import 'prismjs/components/prism-sql.min.js';
import 'prismjs/components/prism-go.min.js';
import 'prismjs/components/prism-python.min.js';
import MyPreset from './mypresset';
import { ConfigService } from './services/config.service';


export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAppInitializer(() => {
      const configService = inject(ConfigService);
      return configService.load();
    }),
    provideHttpClient(withInterceptors([authInterceptor])),
    MessageService,
    provideAnimationsAsync(),
    provideClientHydration(),
    provideMarkdown({
      markedOptions: {
        provide: MARKED_OPTIONS,
        useValue: {
          gfm: true,
          breaks: false,
          pedantic: false,
        }
      },
      clipboardOptions: {
        provide: CLIPBOARD_OPTIONS,
        useValue: {
          buttonComponent : ClipboardButtonComponent,

      },
    }
    }),
    providePrimeNG({
      theme: {
        preset: MyPreset,
        options:{
        }
      }
    })
  ]
};
