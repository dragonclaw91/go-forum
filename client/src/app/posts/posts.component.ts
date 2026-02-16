import { Component } from '@angular/core';
import { MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatSelectModule} from '@angular/material/select';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import { SelectModule } from 'primeng/select';
import { FormsModule } from '@angular/forms'; 
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { CardModule } from 'primeng/card';
import { Timestamp } from 'rxjs';
import { DatePipe } from '@angular/common';
import { DialogModule } from 'primeng/dialog';
import { ButtonModule } from 'primeng/button';
import { TextareaModule } from 'primeng/textarea';

@Component({
	selector: 'app-root',
	imports: [  
    TextareaModule,
    ButtonModule,
    DialogModule,
    DatePipe,
    CardModule,
    InputIconModule,
    IconFieldModule ,
    FormsModule,
 SelectModule ,
    MatSlideToggleModule,
     MatSelectModule, 
     MatFormFieldModule, 
     MatInputModule 
    ],
  templateUrl: './posts.component.html',
  styles:` @use '@angular/material' as mat;
  .my-dropdown-panel {
    @include mat.form-field-overrides((
      outlined-focus-outline-color: orange;
    filled-focus-active-indicator-color: red;
  ));
  }
  `,
  styleUrl: './posts.component.css'
})
export class PostsComponent {
  cities: any[];
  selectedCity: any;

  subposts: any[];

  isVisible = true;

toggleVisibility() {
  this.isVisible = !this.isVisible

}


  constructor() {


    
    this.subposts  = [
      {
        id : 0,
 
      description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
      error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
      nobis,
      culpa ratione quam perferendis esse, cupiditate neque
      quas!`,
   
      created_at : '2025-01-07 10:26:56.904493',
  
    creator_id: 22,
    deleted_post: false,
    votes: 10,
    replies: 99
      },
    {
      id : 1,
 
      description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
      error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
      nobis,
      culpa ratione quam perferendis esse, cupiditate neque
      quas!`,
   
      created_at : '2025-01-07 10:26:56.904493',
  
    creator_id: 22,
    deleted_post: false,
    votes: 30,
    replies:15
    },
    {
      id : 2,
 
      description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
      error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
      nobis,
      culpa ratione quam perferendis esse, cupiditate neque
      quas!`,
   
      created_at : '2025-01-07 10:26:56.904493',
  
    creator_id: 22,
    deleted_post: false,
    votes:56,
    replies:87
    },
    {
      id : 3,
 
      description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
      error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
      nobis,
      culpa ratione quam perferendis esse, cupiditate neque
      quas!`,
   
      created_at : '2025-01-07 10:26:56.904493',
  
    creator_id: 22,
    deleted_post: false,
    votes:0,
    replies: 66
    },
    {
      id : 0,
 
      description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
      error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
      nobis,
      culpa ratione quam perferendis esse, cupiditate neque
      quas!`,
   
      created_at : '2025-01-07 10:26:56.904493',
  
    creator_id: 22,
    deleted_post: false,
    votes:44,
    replies:58
    },
    {
      id : 0,

    description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
    error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
    nobis,
    culpa ratione quam perferendis esse, cupiditate neque
    quas!`,
 
    created_at : '2025-01-07 10:26:56.904493',

  creator_id: 22,
  deleted_post: false,
  votes: 10,
  replies: 99
    },
  {
    id : 1,

    description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
    error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
    nobis,
    culpa ratione quam perferendis esse, cupiditate neque
    quas!`,
 
    created_at : '2025-01-07 10:26:56.904493',

  creator_id: 22,
  deleted_post: false,
  votes: 30,
  replies:15
  },
  {
    id : 2,

    description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
    error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
    nobis,
    culpa ratione quam perferendis esse, cupiditate neque
    quas!`,
 
    created_at : '2025-01-07 10:26:56.904493',

  creator_id: 22,
  deleted_post: false,
  votes:56,
  replies:87
  },
  {
    id : 3,

    description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
    error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
    nobis,
    culpa ratione quam perferendis esse, cupiditate neque
    quas!`,
 
    created_at : '2025-01-07 10:26:56.904493',

  creator_id: 22,
  deleted_post: false,
  votes:0,
  replies: 66
  },
  {
    id : 0,

    description : `Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur
    error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam
    nobis,
    culpa ratione quam perferendis esse, cupiditate neque
    quas!`,
 
    created_at : '2025-01-07 10:26:56.904493',

  creator_id: 22,
  deleted_post: false,
  votes:44,
  replies:58
  }
    ]
    this.cities = [
        {name: 'Likes', code: 'NY'},
        {name: 'Comments', code: 'LDN'},
    ];

}
}
