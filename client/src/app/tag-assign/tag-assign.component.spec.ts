import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TagAssignComponent } from './tag-assign.component';

describe('TagAssignComponent', () => {
  let component: TagAssignComponent;
  let fixture: ComponentFixture<TagAssignComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TagAssignComponent]
    });
    fixture = TestBed.createComponent(TagAssignComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
