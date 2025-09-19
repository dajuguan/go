; llvm_touch_c.ll
define void @llvm_touch_c(i8* %ptr, i8 %val, i64 %n) {
entry:
  %i = alloca i64
  store i64 0, i64* %i
  br label %loop

loop:
  %cur_i = load i64, i64* %i
  %cmp = icmp slt i64 %cur_i, %n
  br i1 %cmp, label %body, label %done

body:
  %elem = getelementptr i8, i8* %ptr, i64 %cur_i
  store i8 %val, i8* %elem
  %next_i = add i64 %cur_i, 1
  store i64 %next_i, i64* %i
  br label %loop

done:
  ret void
}
